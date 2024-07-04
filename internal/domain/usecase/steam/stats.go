package usecasesteam

import (
	"bot/internal/domain/entity"
	"fmt"
	"time"

	"github.com/Newmio/steam_helper"
)

func (s *steam) SearchTradeItems(game string, start, stop int) error {
	ch := make(steam_helper.CursorCh[entity.CheckItem])

	hashNames, err := s.db.GetHashSteamItems(game, int64(start), int64(stop))
	if err != nil {
		return steam_helper.Trace(err)
	}

	links, err := s.db.GetLinkSteamItems(hashNames, game)
	if err != nil {
		return steam_helper.Trace(err)
	}

	go s.r.CheckTradeItem(links, ch)

	for {
		select {

		case item, ok := <-ch:
			if !ok {
				return nil
			}
			if item.Error != nil {
				return steam_helper.Trace(err)
			}

		case <-time.After(time.Minute):
			return steam_helper.Trace(fmt.Errorf("timeout"))
		}
	}
}

func (s *steam) SynchItems(game string) error {
	ch := make(steam_helper.CursorCh[[]entity.SteamItem])

	go s.r.SynchItems(game, ch)

	for {
		select {

		case skin, ok := <-ch:
			if !ok {
				return nil
			}
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
