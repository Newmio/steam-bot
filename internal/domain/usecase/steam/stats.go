package usecasesteam

import (
	"bot/internal/domain/entity"
	"fmt"
	"math"
	"time"

	"github.com/Newmio/steam_helper"
)

func (s *steam) GetHistoryItems(game string, start, stop int) {

}

func (s *steam) CheckTradeItems(game string, start, stop int) error {
	ch := make(steam_helper.CursorCh[entity.CheckItem])

	hashNames, err := s.db.GetHashSteamItems(game, int64(start), int64(stop))
	if err != nil {
		return steam_helper.Trace(err)
	}

	links, err := s.db.GetLinkSteamItems(hashNames, game)
	if err != nil {
		return steam_helper.Trace(err)
	}

	go s.r.CheckTradeItems(links, ch)

	for {
		select {

		case item, ok := <-ch:
			if !ok {
				return nil
			}
			if item.Error != nil {
				return steam_helper.Trace(item.Error)
			}

			maxBuy, minSell := 0.0, math.MaxFloat64

			for cost := range item.Model.Buy {
				costFloat := float64(cost)

				if costFloat > maxBuy {
					maxBuy = costFloat
				}
			}

			for cost := range item.Model.Sell {
				if float64(cost) < minSell {
					minSell = math.Round(float64(cost) * (1 - s.market.Commission/100))
				}
			}

			profit := ((minSell - maxBuy) / maxBuy) * 100

			fmt.Println("-------------------------------------------")
			fmt.Println("hashName", item.Model.HashName)
			fmt.Println("profit", profit)
			fmt.Println("-------------------------------------------")

			if profit > 1.5 { // если процент прибыли больше 1.5%
				history, err := s.db.GetSteamSellHistory(item.Model.HashName, game, 2)
				if err != nil {
					return steam_helper.Trace(err)
				}

				if len(history) == 0 || time.Since(history[len(history)-1].Price.DateTime) >= time.Hour { // если в истории с редиса дата последней покупки старее на 1 час от нынешнего времени
					historyUrl := fmt.Sprintf("https://steamcommunity.com/market/pricehistory/?appid=%s&market_hash_name=%s",
						s.getAppId(game), item.Model.HashName)

					history, err := s.r.GetHistoryItem(historyUrl)
					if err != nil {
						return steam_helper.Trace(err)
					}

					for i := range history {
						history[i].HashName = item.Model.HashName
					}

					if err := s.db.CreateSteamSellHistory(history, game); err != nil {
						return steam_helper.Trace(err)
					}
				}

				var sellsCount int

				for _, value := range history { // считает колво продаж за последние 2 дня
					if time.Since(value.Price.DateTime) <= time.Hour*25*2 {
						sellsCount += value.Price.Count
					}
				}

				if sellsCount > 80 { //если продаж больше 30
					if err := s.db.CreateForSteamTrade(item.Model.HashName, profit); err != nil {
						return steam_helper.Trace(err)
					}
				}
			}

		case <-time.After(time.Minute * 3):
			return steam_helper.Trace(fmt.Errorf("timeout"))
		}
	}
}

func (s *steam) SynchItems(info entity.PaginationInfo[[]entity.SteamItem]) error {
	go s.r.SynchItems(info)

	for {
		select {

		case skin, ok := <-info.Ch:
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

			if err := s.db.CreateHashSteamItems(hashNames, info.Game); err != nil {
				return steam_helper.Trace(err)
			}

			if err := s.db.CreateSteamItems(skin.Model, info.Game); err != nil {
				return steam_helper.Trace(err)
			}

		case <-time.After(time.Minute * 2):
			return steam_helper.Trace(fmt.Errorf("timeout"))
		}
	}
}
