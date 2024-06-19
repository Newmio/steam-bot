package repodb

import (
	reporedis "bot/internal/repository/db/redis"
	reposqlite "bot/internal/repository/db/sqlite"
)

type IDatabase interface{}

type database struct {
	redis  reporedis.IRedis
	sqlite reposqlite.ISqlite
}

func NewDatabase(redis reporedis.IRedis, sqlite reposqlite.ISqlite) IDatabase {
	return &database{redis: redis, sqlite: sqlite}
}
