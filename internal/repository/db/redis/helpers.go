package reporedis

import (
	"context"

	"github.com/Newmio/steam_helper"
	"github.com/redis/go-redis/v9"
)

func (db *redisRepo) GetHelpersForSteamTrade(start, stop int, game string) ([]string, error) {
	var links []string

	result, err := db.db.ZRangeArgsWithScores(context.Background(), redis.ZRangeArgs{
		Key:     game + "_helper_for_steam_trade",
		Start:   start,
		Stop:    stop,
		ByScore: false,
		Rev:     false}).Result()

	if err != nil {
		return nil, steam_helper.Trace(err)
	}

	for _, value := range result {
		links = append(links, value.Member.(string))
	}

	return links, nil
}

func (db *redisRepo) CreateHelpersForSteamTrade(links map[string]float64, game string) error {
	pipe := db.db.TxPipeline()

	for key, value := range links {
		if err := pipe.ZAdd(context.Background(), game+"_helper_for_steam_trade", redis.Z{Score: value, Member: key}).Err(); err != nil {
			return steam_helper.Trace(err)
		}
	}

	if _, err := pipe.Exec(context.Background()); err != nil {
		return steam_helper.Trace(err)
	}

	return nil
}
