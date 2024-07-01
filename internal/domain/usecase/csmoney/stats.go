package usecasecsmoney

import (
	"bot/internal/domain/entity"
	"fmt"
	"time"

	"github.com/Newmio/steam_helper"
)

func (s *csmoney) SynchCSGOSkins(minCost, maxCost float64, minCount int) error {
	ch := make(steam_helper.CursorCh[[]entity.SeleniumSteamSkin])

	go s.r.SynchCsmoneyCSGOSkins(ch)

	for {
		select {

		case skin := <-ch:
			if skin.Error != nil {
				return steam_helper.Trace(skin.Error)
			}

			var skins []entity.SeleniumSteamSkin

			for _, value := range skin.Model {
				normalCost := float64(value.Cost) / 100.0

				if value.Count >= minCount && normalCost >= minCost && normalCost <= maxCost {
					skins = append(skins, value)
				}
			}

			if err := s.db.CreateSeleniumCsmoneySkins(skins); err != nil {
				return steam_helper.Trace(err)
			}

		case <-time.After(5 * time.Minute):
			return steam_helper.Trace(fmt.Errorf("timeout"))
		}
	}
}
