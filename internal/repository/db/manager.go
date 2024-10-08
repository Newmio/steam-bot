package repodb

import (
	"bot/internal/domain/entity"
	reporedis "bot/internal/repository/db/redis"
	reposqlite "bot/internal/repository/db/sqlite"
	"time"
)

//go:generate mockery --name=IDatabase --output=./mocks --case=underscore
type IDatabase interface {
	CreateTables() error
	CreateHashSteamItems(hashNames []string, game string) error
	CreateSteamItems(items []entity.SteamItem, game string) error
	GetHashSteamItems(game string, start, stop int64) ([]string, error)
	GetLinkSteamItems(hashNames []string, game string) ([]string, error)
	// history передавать только по одному предмету
	CreateSteamSellHistory(history []entity.SteamSellHistory, game string) error
	// lastDay - за сколько дней от нынешнего времени выдать историю продаж
	// если lastDay = 0 - то выдать все
	GetSteamSellHistory(hashName, game string, lastDay int) ([]entity.SteamSellHistory, error)
	CreateForSteamTrade(hashName string, profit float64) error
	CreateHelpersForSteamTrade(links map[string]float64, game string) error
	GetHelpersForSteamTrade(start, stop int, game string) ([]string, error)
	CreateItemsRareFloat(items map[string][]entity.FloatItem, game string) error
}

type database struct {
	redis  reporedis.IRedis
	sqlite reposqlite.ISqlite
}

func NewDatabase(redis reporedis.IRedis, sqlite reposqlite.ISqlite) IDatabase {
	return &database{redis: redis, sqlite: sqlite}
}

func (db *database) CreateItemsRareFloat(items map[string][]entity.FloatItem, game string) error {
	return db.redis.CreateItemsRareFloat(items, game)
}

func (db *database) GetHelpersForSteamTrade(start, stop int, game string) ([]string, error) {
	return db.redis.GetHelpersForSteamTrade(start, stop, game)
}

func (db *database) CreateHelpersForSteamTrade(links map[string]float64, game string) error {
	return db.redis.CreateHelpersForSteamTrade(links, game)
}

func (db *database) CreateForSteamTrade(hashName string, profit float64) error {
	return db.redis.CreateForSteamTrade(hashName, profit)
}

func (db *database) GetSteamSellHistory(hashName, game string, lastDay int) ([]entity.SteamSellHistory, error) {
	return db.redis.GetSteamSellHistory(hashName, game, lastDay)
}

func (db *database) CreateSteamSellHistory(history []entity.SteamSellHistory, game string) error {
	var newHist []entity.SteamSellHistory

	for iH := range history {
		if time.Since(history[iH].Price.DateTime) <= time.Hour*24*365 {
			newHist = append(newHist, history[iH])
		}
	}

	if len(newHist) == 0 {
		return nil
	}

	return db.redis.CreateSteamSellHistory(newHist, game)
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
