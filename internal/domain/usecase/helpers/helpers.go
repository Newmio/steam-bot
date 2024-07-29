package usecasehelpers

import (
	repodb "bot/internal/repository/db"
	reposelenium "bot/internal/repository/selenium"

	"github.com/Newmio/steam_helper"
)

type IHelpers interface {
	GetLinksForTradeItem(game string) error
}

type helpers struct {
	r  reposelenium.ISelenium
	db repodb.IDatabase
}

func NewHelpers(r reposelenium.ISelenium, db repodb.IDatabase) IHelpers {
	return &helpers{r: r, db: db}
}

func (s *helpers) GetLinksForTradeItem(game string) error {
	links, err := s.r.GetLinksForTradeItem(game)
	if err != nil {
		return steam_helper.Trace(err)
	}

	if err := s.db.CreateHelpersForSteamTrade(links, game); err != nil {
		return steam_helper.Trace(err)
	}

	return nil
}
