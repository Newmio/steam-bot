package usecasesteam

import (
	"bot/internal/domain/entity"
	"fmt"
	"time"

	"github.com/Newmio/steam_helper"
)

func (s *steam) SearchStickerSkins() error {
	limit, offset := 100, 0

	for {
		skins, err := s.db.GetSeleniumSteamSkins(limit, offset)
		if err != nil{
			return steam_helper.Trace(err)
		}

		for _, value := range skins{
			
		}

		limit += 100
		offset += 100
	}
}

func (s *steam) SynchCSGOSkins(minCost, maxCost float64, minCount int) error {
	ch := make(steam_helper.CursorCh[[]entity.SeleniumSteamSkin])

	go s.r.SynchSteamCSGOSkins(ch)

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

			if err := s.db.CreateSeleniumSteamSkins(skins); err != nil {
				return steam_helper.Trace(err)
			}

		case <-time.After(5 * time.Minute):
			return steam_helper.Trace(fmt.Errorf("timeout"))
		}
	}
}
