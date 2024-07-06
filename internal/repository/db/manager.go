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
	GetHashSteamItems(game string, start, stop int64) ([]string, error)
	GetLinkSteamItems(hashNames []string, game string) ([]string, error)
	CreateSteamSellHistory(history []entity.SteamSellHistory, game string) error
	GetSteamSellHistory(hashName, game string) (entity.SteamSellHistory, error)
	CreateForSteamTrade(hashName string) error
}

type database struct {
	redis  reporedis.IRedis
	sqlite reposqlite.ISqlite
}

func NewDatabase(redis reporedis.IRedis, sqlite reposqlite.ISqlite) IDatabase {
	return &database{redis: redis, sqlite: sqlite}
}

func (db database) CreateForSteamTrade(hashName string) error{
	return db.redis.CreateForSteamTrade(hashName)
}

func (db *database) GetSteamSellHistory(hashName, game string) (entity.SteamSellHistory, error) {
	return db.redis.GetSteamSellHistory(hashName, game)
}

func (db *database) CreateSteamSellHistory(history []entity.SteamSellHistory, game string) error {
	for iH := range history {

		for i, j := 0, len(history[iH].Prices)-1; i < j; i, j = i+1, j-1 {
			history[iH].Prices[i], history[iH].Prices[j] = history[iH].Prices[j], history[iH].Prices[i]
		}
	}
	return db.redis.CreateSteamSellHistory(history, game)
}

func (db *database) GetLinkSteamItems(hashNames []string, game string) ([]string, error) {
	return db.redis.GetLinkSteamItems(hashNames, game)
}

func (db *database) GetHashSteamItems(game string, start, stop int64) ([]string, error) {
	return db.redis.GetHashSteamItems(game, start, stop)
}

func (db *database) CreateHashSteamItems(hashNames []string, game string) error {
	return db.redis.CreateHashSteamItems(hashNames, game)
}

func (db *database) CreateSteamItems(items []entity.SteamItem, game string) error {
	return db.redis.CreateSteamItems(items, game)
}

func (db *database) CreateTables() error {
	return db.sqlite.CreateTables()
}
