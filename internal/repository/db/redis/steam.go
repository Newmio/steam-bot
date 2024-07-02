package reporedis

import (
	"bot/internal/domain/entity"
	"context"
	"encoding/json"
	"fmt"

	"github.com/Newmio/steam_helper"
)

func (db *redisRepo) CreateSteamItems(items []entity.SteamItem, game string) error {
	pipe := db.db.TxPipeline()

	for _, value := range items {
		body, err := json.Marshal(value)
		if err != nil {
			return steam_helper.Trace(err, value)
		}

		if err := pipe.HMSet(context.Background(), fmt.Sprintf("steam_%s_items_%s", game, value.HashName), body).Err(); err != nil {
			return steam_helper.Trace(err)
		}
	}

	if _, err := pipe.Exec(context.Background()); err != nil {
		return steam_helper.Trace(err)
	}

	return nil
}

func (db *redisRepo) CreateHashSteamItems(hashNames []string, game string) error {
	pipe := db.db.TxPipeline()

	for _, value := range hashNames {
		if err := pipe.SAdd(context.Background(), fmt.Sprintf("steam_%s_items", game), value).Err(); err != nil {
			return steam_helper.Trace(err)
		}
	}

	if _, err := pipe.Exec(context.Background()); err != nil {
		return steam_helper.Trace(err)
	}

	return nil
}