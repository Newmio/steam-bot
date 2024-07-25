package usecasesteam

import (
	"bot/internal/domain/entity"
	repodb "bot/internal/repository/db"
	reposelenium "bot/internal/repository/selenium"
)

type ISteam interface {
	SteamAuth() error
	SynchItems(info entity.PaginationInfo[[]entity.SteamItem]) error
	Ping(url string) (string, error)
	CheckItems(game string, start, stop int) error
}

type steam struct {
	r      reposelenium.ISelenium
	db     repodb.IDatabase
	market entity.Market
}

func NewSteam(r reposelenium.ISelenium, db repodb.IDatabase, market entity.Market) ISteam {
	return &steam{r: r, db: db, market: market}
}

func (s *steam) Ping(url string) (string, error) {
	return s.r.Ping(url)
}

func (s *steam) getAppId(game string)string{
	switch game {
	case "csgo":
		return "730"

	default:
		return "error"
	}
}