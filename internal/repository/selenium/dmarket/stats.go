package repodmarket

import (
	"bot/internal/domain/entity"

	"github.com/Newmio/steam_helper"
	"github.com/tebeka/selenium"
)

func (r *dmarket) SynchCSGOSkins(wd selenium.WebDriver, ch steam_helper.CursorCh[[]entity.SeleniumSteamSkin]) error {
	return nil
}
