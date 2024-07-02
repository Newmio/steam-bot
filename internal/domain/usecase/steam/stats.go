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
		// skins, err := s.db.GetSeleniumSteamSkins(limit, offset)
		// if err != nil{
		// 	return steam_helper.Trace(err)
		// }

		// for _, value := range skins{

		// }

		limit += 100
		offset += 100
	}
}

func (s *steam) SynchCSGOItems(game string) error {
	ch := make(steam_helper.CursorCh[[]entity.SteamItem])

	go s.r.SynchCSGOItems(ch)

	for {
		select {

		case skin := <-ch:
			if skin.Error != nil {
				return steam_helper.Trace(skin.Error)
			}

			var hashNames []string

			for _, value := range skin.Model {
				hashNames = append(hashNames, value.HashName)
			}

			if err := s.db.CreateHashSteamItems(hashNames, game); err != nil {
				return steam_helper.Trace(err)
			}

		case <-time.After(time.Minute):
			return steam_helper.Trace(fmt.Errorf("timeout"))
		}
	}
}
