package reporedis

import (
	"bot/internal/domain/entity"

	"github.com/redis/go-redis/v9"
)

type IRedis interface {
	CreateHashSteamItems(hashNames []string, game string) error
	CreateSteamItems(items []entity.SteamItem, game string) error
}

type redisRepo struct {
	db *redis.Client
}

func NewRedis(db *redis.Client) IRedis {
	return &redisRepo{db: db}
}
