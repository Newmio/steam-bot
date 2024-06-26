package reporedis

import (
	"bot/internal/domain/entity"

	"github.com/redis/go-redis/v9"
)

type IRedis interface {
	CreateSeleniumSteamSkins(skins []entity.SeleniumSteamSkin) error
	CreateSteamSkins(skins []entity.DbSteamSkins) error
	CreateStickerSkins(skins []entity.DbSteamSkins) error
	CreateFloatSkins(skins []entity.DbSteamSkins) error
	CreatePatternSkins(skins []entity.DbSteamSkins) error
	CreateBetweenSkins(skins []entity.DbSteamSkins) error
	GetSeleniumSteamSkins(start, stop int) ([]entity.SeleniumSteamSkin, error)
	GetSteamSkins(start, stop int) ([]entity.DbSteamSkins, error)
	GetStickerSkins(start, stop int) ([]entity.DbSteamSkins, error)
	GetFloatSkins(start, stop int) ([]entity.DbSteamSkins, error)
	GetPatternSkins(start, stop int) ([]entity.DbSteamSkins, error)
	GetBetweenSkins(start, stop int) ([]entity.DbSteamSkins, error)
	CreateSeleniumDmarketSkins(skins []entity.SeleniumSteamSkin) error
}

type redisRepo struct {
	db *redis.Client
}

func NewRedis(db *redis.Client) IRedis {
	return &redisRepo{db: db}
}
