package reposteam

import (
	"bot/internal/domain/entity"

	"github.com/Newmio/steam_helper"
	"github.com/tebeka/selenium"
)

type ISteam interface {
	Login(wd selenium.WebDriver, user entity.SteamUser) (string, error)
	SynchCSGOItems(wd selenium.WebDriver, ch steam_helper.CursorCh[[]entity.SteamItem])
}

type steam struct{}

func NewSteam() ISteam {
	return &steam{}
}
