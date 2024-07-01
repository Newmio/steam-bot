package reposteam

import (
	"bot/internal/domain/entity"

	"github.com/Newmio/steam_helper"
	"github.com/tebeka/selenium"
)

type ISteam interface {
	Login(wd selenium.WebDriver, user entity.SteamUser) (string, error)
	SynchCSGOSkins(wd selenium.WebDriver, ch steam_helper.CursorCh[[]entity.SeleniumSteamSkin])
}

type steam struct{}

func NewSteam() ISteam {
	return &steam{}
}
