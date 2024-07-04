package reporedis

import (
	"bot/internal/domain/entity"
	"context"
	"fmt"

	"github.com/Newmio/steam_helper"
	"github.com/redis/go-redis/v9"
)

func (db *redisRepo) GetLinkSteamItems(hashNames []string, game string) ([]string, error) {
	var links []string
	var stringCmd []*redis.StringCmd
	pipe := db.db.TxPipeline()

	for _, value := range hashNames {
		stringCmd = append(stringCmd, pipe.HGet(context.Background(), fmt.Sprintf("steam_%s_items_%s", game, value), "link"))
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
		if err := pipe.HSet(context.Background(), fmt.Sprintf("steam_%s_items_%s", game, value.HashName), value).Err(); err != nil {
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
