package reporedis

import (
	"bot/internal/domain/entity"
	"context"
	"fmt"

	"github.com/Newmio/steam_helper"
	"github.com/redis/go-redis/v9"
)

func (db *redisRepo) CreateForSteamTrade(hashName string) error {
	return db.db.ZAdd(context.Background(), "for_trade_steam", redis.Z{Score: 0, Member: hashName}).Err()
}

func (db *redisRepo) GetSteamSellHistory(hashName, game string) (entity.SteamSellHistory, error) {
	var history entity.SteamSellHistory
	var respModel map[string]interface{}
	var respPrices map[string]interface{}

	if err := db.db.HGetAll(context.Background(), fmt.Sprintf("steam_%s_sell_history:%s", game, hashName)).Scan(&respModel); err != nil {
		return history, steam_helper.Trace(err)
	}

	for i := range int(respModel["countPrices"].(float64)) {
		if err := db.db.HGetAll(context.Background(), fmt.Sprintf("steam_%s_sell_history:%s[%d]", game, hashName, i)).Scan(&respPrices); err != nil {
			return history, steam_helper.Trace(err)
		}

		date, err := steam_helper.TimeParse(respPrices["DateTime"].(string))
		if err != nil {
			return history, steam_helper.Trace(err)
		}

		respPrices["DateTime"] = date
	}

	respModel["Prices"] = respPrices

	if err := steam_helper.MapToStruct(respModel, &history); err != nil {
		return history, steam_helper.Trace(err)
	}

	return history, nil
}

func (db *redisRepo) CreateSteamSellHistory(history []entity.SteamSellHistory, game string) error {
	pipe := db.db.TxPipeline()

	for _, value := range history {
		model, err := steam_helper.StructToMap(value)
		if err != nil {
			return steam_helper.Trace(err)
		}

		for i := len(value.Prices) - 1; i >= 0; i-- {
			arrModel, err := steam_helper.StructToMap(value.Prices[i])
			if err != nil {
				return steam_helper.Trace(err)
			}

			arrModel["DateTime"] = steam_helper.TimeFormat(value.Prices[i].DateTime)

			if err := pipe.HSet(context.Background(), fmt.Sprintf("steam_%s_sell_history:%s[%d]", game, value.HashName, i), arrModel).Err(); err != nil {
				return steam_helper.Trace(err)
			}
		}

		model["countPrices"] = len(value.Prices)

		if err := pipe.HSet(context.Background(), fmt.Sprintf("steam_%s_sell_history:%s", game, value.HashName), model).Err(); err != nil {
			return steam_helper.Trace(err)
		}
	}

	if _, err := pipe.Exec(context.Background()); err != nil {
		return steam_helper.Trace(err)
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
