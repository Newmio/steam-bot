package repocsmoney

import (
	"bot/internal/domain/entity"

	"github.com/Newmio/steam_helper"
	"github.com/tebeka/selenium"
)

type ICsmoney interface {
	SynchCSGOSkins(wd selenium.WebDriver, ch steam_helper.CursorCh[[]entity.SeleniumSteamSkin]) error
}

type csmoney struct {
}

func NewCsmoney() ICsmoney {
	return &csmoney{}
}
