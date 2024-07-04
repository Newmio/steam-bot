package reposteam

import (
	"bot/internal/domain/entity"

	"github.com/Newmio/steam_helper"
	"github.com/tebeka/selenium"
)

type ISteam interface {
	Login(wd selenium.WebDriver, user entity.SteamUser) (string, error)
	SynchItems(wd selenium.WebDriver, game string, ch steam_helper.CursorCh[[]entity.SteamItem])
	CheckTradeItems(wd selenium.WebDriver, links []string, ch steam_helper.CursorCh[entity.CheckItem])
}

type steam struct{}

func NewSteam() ISteam {
	return &steam{}
}
