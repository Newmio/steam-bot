package repodb

import (
	"bot/internal/domain/entity"
	reporedis "bot/internal/repository/db/redis"
	reposqlite "bot/internal/repository/db/sqlite"

	"github.com/Newmio/steam_helper"
)

type IDatabase interface {
	CreateTables() error
	CreateSteamSkins(skins []entity.DbSteamSkins) error
	CreateStickerSkins(skins []entity.DbSteamSkins) error
	CreateFloatSkins(skins []entity.DbSteamSkins) error
	CreatePatternSkins(skins []entity.DbSteamSkins) error
	GetSteamSkins(limit, offset int) ([]entity.DbSteamSkins, error)
	GetStickerSkins(limit, offset int) ([]entity.DbSteamSkins, error)
	GetFloatSkins(limit, offset int) ([]entity.DbSteamSkins, error)
	GetPatternSkins(limit, offset int) ([]entity.DbSteamSkins, error)
	CreateSeleniumSteamSkins(skins []entity.SeleniumSteamSkin) error
	CreateSeleniumDmarketSkins(skins []entity.SeleniumSteamSkin) error
	CreateSeleniumCsmoneySkins(skins []entity.SeleniumSteamSkin) error
}

type database struct {
	redis  reporedis.IRedis
	sqlite reposqlite.ISqlite
}

func NewDatabase(redis reporedis.IRedis, sqlite reposqlite.ISqlite) IDatabase {
	return &database{redis: redis, sqlite: sqlite}
}

func (db *database) CreateSeleniumCsmoneySkins(skins []entity.SeleniumSteamSkin) error{
	return db.redis.CreateSeleniumCsmoneySkins(skins)
}

func (db *database) CreateSeleniumDmarketSkins(skins []entity.SeleniumSteamSkin) error{
	return db.redis.CreateSeleniumDmarketSkins(skins)
}

func (db *database) CreateBetweenSkins(skins []entity.DbSteamSkins) error {
	return db.redis.CreateBetweenSkins(skins)
}

func (db *database) CreateSeleniumSteamSkins(skins []entity.SeleniumSteamSkin) error {
	return db.redis.CreateSeleniumSteamSkins(skins)
}

func (db *database) CreateSteamSkins(skins []entity.DbSteamSkins) error {
	return db.redis.CreateSteamSkins(skins)
}

func (db *database) CreateStickerSkins(skins []entity.DbSteamSkins) error {
	return db.redis.CreateStickerSkins(skins)
}

func (db *database) CreateFloatSkins(skins []entity.DbSteamSkins) error {
	return db.redis.CreateFloatSkins(skins)
}

func (db *database) CreatePatternSkins(skins []entity.DbSteamSkins) error {
	return db.redis.CreatePatternSkins(skins)
}

func (db *database) GetSteamSkins(limit, offset int) ([]entity.DbSteamSkins, error) {
	skins, err := db.redis.GetSteamSkins(offset, limit+offset-1)
	if err != nil {
		return nil, steam_helper.Trace(err)
	}

	if len(skins) > 0 {
		return skins, nil
	}

	return db.sqlite.GetSteamSkins(limit, offset)
}

func (db *database) GetStickerSkins(limit, offset int) ([]entity.DbSteamSkins, error) {
	skins, err := db.redis.GetStickerSkins(offset, limit+offset-1)
	if err != nil {
		return nil, steam_helper.Trace(err)
	}

	if len(skins) > 0 {
		return skins, nil
	}

	return db.sqlite.GetStickerSkins(limit, offset)
}

func (db *database) GetFloatSkins(limit, offset int) ([]entity.DbSteamSkins, error) {
	skins, err := db.redis.GetFloatSkins(offset, limit+offset-1)
	if err != nil {
		return nil, steam_helper.Trace(err)
	}

	if len(skins) > 0 {
		return skins, nil
	}

	return db.sqlite.GetFloatSkins(limit, offset)
}

func (db *database) GetPatternSkins(limit, offset int) ([]entity.DbSteamSkins, error) {
	skins, err := db.redis.GetPatternSkins(offset, limit+offset-1)
	if err != nil {
		return nil, steam_helper.Trace(err)
	}

	if len(skins) > 0 {
		return skins, nil
	}

	return db.sqlite.GetPatternSkins(limit, offset)
}

func (db *database) GetBetweenSkins(limit, offset int) ([]entity.DbSteamSkins, error) {
	skins, err := db.redis.GetBetweenSkins(offset, limit+offset-1)
	if err != nil {
		return nil, steam_helper.Trace(err)
	}

	if len(skins) > 0 {
		return skins, nil
	}

	return db.sqlite.GetBetweenSkins(limit, offset)
}

func (db *database) CreateTables() error {
	return db.sqlite.CreateTables()
}
