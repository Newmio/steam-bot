package reporedis

import (
	"bot/internal/domain/entity"
	"context"
	"encoding/json"
	"fmt"

	"github.com/Newmio/steam_helper"
)

func (db *redisRepo) CreateItemsRareFloat(items map[string][]entity.FloatItem, game string) error {
	pipe := db.db.TxPipeline()

	for key, values := range items {

		for _, value := range values {
			body, err := json.Marshal(value)
			if err != nil {
				return steam_helper.Trace(err)
			}

			if err := pipe.RPush(context.Background(), fmt.Sprintf("rare_float_%s_items:%s", game, key), body).Err(); err != nil {
				return steam_helper.Trace(err)
			}
		}
	}

	if _, err := pipe.Exec(context.Background()); err != nil {
		return steam_helper.Trace(err)
	}

	return nil
}
