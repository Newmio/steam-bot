package reporedis

import (
	"bot/internal/domain/entity"
	"context"
	"encoding/json"

	"github.com/Newmio/steam_helper"
)

func (db *redisRepo) CreateSeleniumDmarketSkins(skins []entity.SeleniumSteamSkin) error {
	for _, value := range skins {

		body, err := json.Marshal(value)
		if err != nil {
			return steam_helper.Trace(err, value)
		}

		err = db.db.RPush(context.Background(), "selenium_dmarket_skins", string(body)).Err()
		if err != nil {
			return steam_helper.Trace(err)
		}
	}

	return nil
}
