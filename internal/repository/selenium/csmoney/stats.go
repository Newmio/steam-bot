package repocsmoney

import (
	"bot/internal/domain/entity"

	"github.com/Newmio/steam_helper"
	"github.com/tebeka/selenium"
)

func (r *csmoney) SynchCSGOSkins(wd selenium.WebDriver, ch steam_helper.CursorCh[[]entity.SeleniumSteamSkin]) error {
	return nil
}
