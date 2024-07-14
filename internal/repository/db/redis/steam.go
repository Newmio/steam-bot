package reporedis

import (
	"bot/internal/domain/entity"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Newmio/steam_helper"
	"github.com/redis/go-redis/v9"
)

func (db *redisRepo) CreateForSteamTrade(hashName string) error {
	return db.db.ZAdd(context.Background(), "for_trade_steam", redis.Z{Score: 0, Member: hashName}).Err()
}

// lastDay - за сколько дней от нынешнего времени выдать историю продаж
func (db *redisRepo) GetSteamSellHistory(hashName, game string, lastDay int) ([]entity.SteamSellHistory, error) {
	var history []entity.SteamSellHistory
	t := time.Now()

	start, stop := 0, 100
	for {
		resp, err := db.db.LRange(context.Background(), fmt.Sprintf("steam_%s_sell_history:%s", game, hashName), int64(start), int64(stop)).Result()
		if err != nil {
			return history, steam_helper.Trace(err)
		}

		if len(resp) == 0 {
			return history, nil
		}

		for _, value := range resp {
			var historyItem entity.SteamSellHistory

			if err := json.Unmarshal([]byte(value), &historyItem); err != nil {
				return history, steam_helper.Trace(err)
			}

			if t.Sub(historyItem.Price.DateTime) <= time.Hour*24*time.Duration(lastDay) || lastDay == 0 {
				history = append(history, historyItem)
			} else {
				return history, nil
			}
		}

		start += 100
		stop += 100
	}
}

func (db *redisRepo) CreateSteamSellHistory(history []entity.SteamSellHistory, game string) error {
	index := 0
	pipe := db.db.TxPipeline()
	hash := fmt.Sprintf("steam_%s_sell_history:%s", game, history[len(history)-1].HashName)

	respCheck, err := db.db.LIndex(context.Background(), hash, -1).Result()
	if err != nil {
		if err != redis.Nil {
			return steam_helper.Trace(err)
		}
	}

	if respCheck != "" {
		var historyItem entity.SteamSellHistory

		if err := json.Unmarshal([]byte(respCheck), &historyItem); err != nil {
			return steam_helper.Trace(err)
		}

		fmt.Println("---- 1 len", len(history))

		for i := len(history) - 1; i >= 0; i-- {
			if !history[i].Price.DateTime.After(historyItem.Price.DateTime) {
				history = history[i+1:]
				break
			}
		}

		fmt.Println("---- 2 len", len(history))
	}

	for _, value := range history {

		body, err := json.Marshal(value)
		if err != nil {
			return steam_helper.Trace(err, value)
		}

		if err := pipe.RPush(context.Background(), hash, body).Err(); err != nil {
			return steam_helper.Trace(err)
		}

		index++

		if index == 100 {
			if _, err := pipe.Exec(context.Background()); err != nil {
				return steam_helper.Trace(err)
			}
			pipe.Discard()
		}
	}

	if index != 0 {
		if _, err := pipe.Exec(context.Background()); err != nil {
			return steam_helper.Trace(err)
		}
	}

	return nil
}

func (db *redisRepo) GetLinkSteamItems(hashNames []string, game string) ([]string, error) {
	var links []string
	var stringCmd []*redis.StringCmd
	pipe := db.db.TxPipeline()

	for _, value := range hashNames {
		stringCmd = append(stringCmd, pipe.HGet(context.Background(), fmt.Sprintf("steam_%s_item:%s", game, value), "Link"))
	}

	if _, err := pipe.Exec(context.Background()); err != nil {
		return nil, steam_helper.Trace(err)
	}

	for _, value := range stringCmd {
		resp, err := value.Result()
		if err != nil {
			return nil, steam_helper.Trace(err)
		}
		links = append(links, resp)
	}

	return links, nil
}

func (db *redisRepo) CreateSteamItems(items []entity.SteamItem, game string) error {
	pipe := db.db.TxPipeline()

	for _, value := range items {
		model, err := steam_helper.StructToMap(value)
		if err != nil {
			return steam_helper.Trace(err)
		}

		if err := pipe.HSet(context.Background(), fmt.Sprintf("steam_%s_item:%s", game, value.HashName), model).Err(); err != nil {
			return steam_helper.Trace(err)
		}
	}

	if _, err := pipe.Exec(context.Background()); err != nil {
		return steam_helper.Trace(err)
	}

	return nil
}

func (db *redisRepo) GetHashSteamItems(game string, start, stop int64) ([]string, error) {
	return db.db.ZRange(context.Background(), fmt.Sprintf("steam_%s_items", game), start, stop).Result()
}

func (db *redisRepo) CreateHashSteamItems(hashNames []string, game string) error {
	pipe := db.db.TxPipeline()

	for _, value := range hashNames {
		if err := pipe.ZAdd(context.Background(), fmt.Sprintf("steam_%s_items", game), redis.Z{Score: 0, Member: value}).Err(); err != nil {
			return steam_helper.Trace(err)
		}
	}

	if _, err := pipe.Exec(context.Background()); err != nil {
		return steam_helper.Trace(err)
	}

	return nil
}
