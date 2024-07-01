package reposteam

import (
	"bot/internal/domain/entity"
	"context"
	"strconv"

	"github.com/Newmio/steam_helper"

	"github.com/tebeka/selenium"
)

func (r *steam) SynchCSGOSkins(wd selenium.WebDriver, ch steam_helper.CursorCh[[]entity.SeleniumSteamSkin]) {
	if err := wd.Get("https://steamcommunity.com/market/search?appid=730#p1_popular_desc"); err != nil {
		ch.WriteError(context.Background(), steam_helper.Trace(err, wd))
		return
	}

	steam_helper.SleepRandom(9000, 10000)

	start, err := steam_helper.GetStartMousePosition(wd)
	if err != nil {
		ch.WriteError(context.Background(), steam_helper.Trace(err, wd))
		return
	}

	for {
		var steamSkins []entity.SeleniumSteamSkin

		skins, err := wd.FindElements(selenium.ByCSSSelector, ".market_listing_row_link")
		if err != nil {
			ch.WriteError(context.Background(), steam_helper.Trace(err, wd))
			return
		}

		for _, skin := range skins {

			hashNameElement, err := skin.FindElement(selenium.ByCSSSelector, ".market_listing_row.market_recent_listing_row.market_listing_searchresult")
			if err != nil {
				ch.WriteError(context.Background(), steam_helper.Trace(err, skin))
				return
			}

			hashName, err := hashNameElement.GetAttribute("data-hash-name")
			if err != nil {
				ch.WriteError(context.Background(), steam_helper.Trace(err, hashNameElement))
				return
			}

			ruNameElement, err := skin.FindElement(selenium.ByCSSSelector, ".market_listing_item_name")
			if err != nil {
				ch.WriteError(context.Background(), steam_helper.Trace(err, skin))
				return
			}

			ruName, err := ruNameElement.Text()
			if err != nil {
				ch.WriteError(context.Background(), steam_helper.Trace(err, ruNameElement))
				return
			}

			costMainElement, err := skin.FindElement(selenium.ByCSSSelector, ".market_table_value.normal_price")
			if err != nil {
				ch.WriteError(context.Background(), steam_helper.Trace(err, skin))
				return
			}

			costElement, err := costMainElement.FindElement(selenium.ByCSSSelector, ".normal_price")
			if err != nil {
				ch.WriteError(context.Background(), steam_helper.Trace(err, costMainElement))
				return
			}

			costStr, err := costElement.GetAttribute("data-price")
			if err != nil {
				ch.WriteError(context.Background(), steam_helper.Trace(err, costElement))
				return
			}

			countElement, err := skin.FindElement(selenium.ByCSSSelector, ".market_listing_num_listings_qty")
			if err != nil {
				ch.WriteError(context.Background(), steam_helper.Trace(err, skin))
				return
			}

			countStr, err := countElement.Text()
			if err != nil {
				ch.WriteError(context.Background(), steam_helper.Trace(err, countElement))
				return
			}

			link, err := skin.GetAttribute("href")
			if err != nil{
				ch.WriteError(context.Background(), steam_helper.Trace(err, skin))
				return
			}

			nextBtn, err := wd.FindElement(selenium.ByCSSSelector, ".pagebtn")
			if err != nil {
				ch.WriteError(context.Background(), steam_helper.Trace(err, wd))
				return
			}

			end, err := steam_helper.MoveMouseAndClick(wd, nextBtn, start)
			if err != nil {
				ch.WriteError(context.Background(), steam_helper.Trace(err, nextBtn))
				return
			}
			start = end

			cost, err := strconv.Atoi(costStr)
			if err != nil {
				ch.WriteError(context.Background(), steam_helper.Trace(err))
				return
			}

			count, err := strconv.Atoi(countStr)
			if err != nil {
				ch.WriteError(context.Background(), steam_helper.Trace(err))
				return
			}

			steamSkins = append(steamSkins, entity.SeleniumSteamSkin{
				HashName: hashName,
				RuName:   ruName,
				Cost:     cost,
				Count:    count,
				Link: link,
			})
		}

		ch.WriteModel(context.Background(), steamSkins)
	}
}
