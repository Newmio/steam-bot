package reporedis

import (
	"bot/internal/domain/entity"
	"context"
	"encoding/json"

	"github.com/Newmio/steam_helper"
)

func (db *redisRepo) GetPatternSkins(start, stop int) ([]entity.DbSteamSkins, error) {
	var skins []entity.DbSteamSkins

	c, err := db.db.LRange(context.Background(), "_skins", int64(start), int64(stop)).Result()
	if err != nil {
		return nil, steam_helper.Trace(err)
	}

	for _, value := range c {
		var skin entity.DbSteamSkins

		if err := json.Unmarshal([]byte(value), &skin); err != nil {
			return nil, steam_helper.Trace(err)
		}

		skins = append(skins, skin)
	}

	return skins, nil
}

func (db *redisRepo) GetFloatSkins(start, stop int) ([]entity.DbSteamSkins, error) {
	var skins []entity.DbSteamSkins

	c, err := db.db.LRange(context.Background(), "_skins", int64(start), int64(stop)).Result()
	if err != nil {
		return nil, steam_helper.Trace(err)
	}

	for _, value := range c {
		var skin entity.DbSteamSkins

		if err := json.Unmarshal([]byte(value), &skin); err != nil {
			return nil, steam_helper.Trace(err)
		}

		skins = append(skins, skin)
	}

	return skins, nil
}

func (db *redisRepo) GetStickerSkins(start, stop int) ([]entity.DbSteamSkins, error) {
	var skins []entity.DbSteamSkins

	c, err := db.db.LRange(context.Background(), "_skins", int64(start), int64(stop)).Result()
	if err != nil {
		return nil, steam_helper.Trace(err)
	}

	for _, value := range c {
		var skin entity.DbSteamSkins

		if err := json.Unmarshal([]byte(value), &skin); err != nil {
			return nil, steam_helper.Trace(err)
		}

		skins = append(skins, skin)
	}

	return skins, nil
}

func (db *redisRepo) GetSteamSkins(start, stop int) ([]entity.DbSteamSkins, error) {
	var skins []entity.DbSteamSkins

	c, err := db.db.LRange(context.Background(), "_skins", int64(start), int64(stop)).Result()
	if err != nil {
		return nil, steam_helper.Trace(err)
	}

	for _, value := range c {
		var skin entity.DbSteamSkins

		if err := json.Unmarshal([]byte(value), &skin); err != nil {
			return nil, steam_helper.Trace(err)
		}

		skins = append(skins, skin)
	}

	return skins, nil
}

func (db *redisRepo) GetSeleniumSteamSkins(start, stop int) ([]entity.SeleniumSteamSkin, error) {
	var skins []entity.SeleniumSteamSkin

	c, err := db.db.LRange(context.Background(), "_skins", int64(start), int64(stop)).Result()
	if err != nil {
		return nil, steam_helper.Trace(err)
	}

	for _, value := range c {
		var skin entity.SeleniumSteamSkin

		if err := json.Unmarshal([]byte(value), &skin); err != nil {
			return nil, steam_helper.Trace(err)
		}

		skins = append(skins, skin)
	}

	return skins, nil
}

func (db *redisRepo) CreatePatternSkins(skins []entity.DbSteamSkins) error {
	for _, value := range skins {

		body, err := json.Marshal(value)
		if err != nil {
			return steam_helper.Trace(err, value)
		}

		err = db.db.RPush(context.Background(), "pattern_skins", string(body)).Err()
		if err != nil {
			steam_helper.Trace(err)
		}
	}

	return nil
}

func (db *redisRepo) CreateFloatSkins(skins []entity.DbSteamSkins) error {
	for _, value := range skins {

		body, err := json.Marshal(value)
		if err != nil {
			return steam_helper.Trace(err, value)
		}

		err = db.db.RPush(context.Background(), "float_skins", string(body)).Err()
		if err != nil {
			steam_helper.Trace(err)
		}
	}

	return nil
}

func (db *redisRepo) CreateStickerSkins(skins []entity.DbSteamSkins) error {
	for _, value := range skins {

		body, err := json.Marshal(value)
		if err != nil {
			return steam_helper.Trace(err, value)
		}

		err = db.db.RPush(context.Background(), "sticker_skins", string(body)).Err()
		if err != nil {
			steam_helper.Trace(err)
		}
	}

	return nil
}

func (db *redisRepo) CreateSteamSkins(skins []entity.DbSteamSkins) error {
	for _, value := range skins {

		body, err := json.Marshal(value)
		if err != nil {
			return steam_helper.Trace(err, value)
		}

		err = db.db.RPush(context.Background(), "steam_skins", string(body)).Err()
		if err != nil {
			steam_helper.Trace(err)
		}
	}

	return nil
}

func (db *redisRepo) CreateSeleniumSteamSkins(skins []entity.SeleniumSteamSkin) error {
	for _, value := range skins {

		body, err := json.Marshal(value)
		if err != nil {
			return steam_helper.Trace(err, value)
		}

		err = db.db.RPush(context.Background(), "selenium_steam_skins", string(body)).Err()
		if err != nil {
			return steam_helper.Trace(err)
		}
	}

	return nil
}
