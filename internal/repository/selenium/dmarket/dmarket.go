package repodmarket

import (
	"bot/internal/domain/entity"

	"github.com/Newmio/steam_helper"
	"github.com/tebeka/selenium"
)

type IDmarket interface {
	SynchCSGOSkins(wd selenium.WebDriver, ch steam_helper.CursorCh[[]entity.SeleniumSteamSkin]) error
}

type dmarket struct {
}

func NewDmarket() IDmarket {
	return &dmarket{}
}
