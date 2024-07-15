package reporedis

import (
	"bot/internal/domain/entity"

	"github.com/redis/go-redis/v9"
)

type IRedis interface {
	CreateHashSteamItems(hashNames []string, game string) error
	CreateSteamItems(items []entity.SteamItem, game string) error
	GetHashSteamItems(game string, start, stop int64) ([]string, error)
	GetLinkSteamItems(hashNames []string, game string) ([]string, error)
	CreateSteamSellHistory(history []entity.SteamSellHistory, game string) error
	GetSteamSellHistory(hashName, game string, lastDay int) ([]entity.SteamSellHistory, error)
	CreateForSteamTrade(hashName string, profit float64) error
}

type redisRepo struct {
	db *redis.ClusterClient
}

func NewRedis(db *redis.ClusterClient) IRedis {
	return &redisRepo{db: db}
}
