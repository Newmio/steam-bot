package repodb

import (
	"bot/internal/domain/entity"
	reporedis "bot/internal/repository/db/redis"
	reposqlite "bot/internal/repository/db/sqlite"
)

type IDatabase interface {
	CreateTables() error
	CreateHashSteamItems(hashNames []string, game string) error
	CreateSteamItems(items []entity.SteamItem, game string) error
}

type database struct {
	redis  reporedis.IRedis
	sqlite reposqlite.ISqlite
}

func NewDatabase(redis reporedis.IRedis, sqlite reposqlite.ISqlite) IDatabase {
	return &database{redis: redis, sqlite: sqlite}
}
func (db *database) CreateTables() error {
	return db.sqlite.CreateTables()
}

func (db *database) CreateHashSteamItems(hashNames []string, game string) error {
	return db.redis.CreateHashSteamItems(hashNames, game)
}

func (db *database) CreateSteamItems(items []entity.SteamItem, game string) error {
	return db.redis.CreateSteamItems(items, game)
}
