package reporedis

import "github.com/redis/go-redis/v9"

type IRedis interface {
}

type redisRepo struct {
	db *redis.Client
}

func NewRedis(db *redis.Client) IRedis {
	return &redisRepo{db: db}
}
