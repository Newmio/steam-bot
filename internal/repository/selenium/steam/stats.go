package reposteam

import (
	"bot/internal/domain/entity"
	"context"
	"fmt"
	"strconv"

	"github.com/Newmio/steam_helper"

	"github.com/tebeka/selenium"
)

func (r *steam) GetCSGOSkins(wd selenium.WebDriver, ch steam_helper.CursorCh[[]entity.SeleniumSteamSkin]) error {
	if err := wd.Get("https://steamcommunity.com/market/search?appid=730#p1_popular_desc"); err != nil {
		return steam_helper.Trace(err, wd)
	}

	steam_helper.SleepRandom(9000, 10000)

	start, err := steam_helper.GetStartMousePosition(wd)
	if err != nil {
		return steam_helper.Trace(err, wd)
	}

	for {
		var steamSkins []entity.SeleniumSteamSkin

		skins, err := wd.FindElements(selenium.ByCSSSelector, ".market_listing_row_link")
		if err != nil {
			return steam_helper.Trace(err, wd)
		}

		for _, skin := range skins {

			hashNameElement, err := skin.FindElement(selenium.ByCSSSelector, ".market_listing_row.market_recent_listing_row.market_listing_searchresult")
			if err != nil {
				return steam_helper.Trace(err, skin)
			}

			hashName, err := hashNameElement.GetAttribute("data-hash-name")
			if err != nil {
				return steam_helper.Trace(err, hashNameElement)
			}

			ruNameElement, err := skin.FindElement(selenium.ByCSSSelector, ".market_listing_item_name")
			if err != nil {
				return steam_helper.Trace(err, skin)
			}

			ruName, err := ruNameElement.Text()
			if err != nil {
				return steam_helper.Trace(err, ruNameElement)
			}

			costMainElement, err := skin.FindElement(selenium.ByCSSSelector, ".market_table_value.normal_price")
			if err != nil {
				return steam_helper.Trace(err, skin)
			}

			costElement, err := costMainElement.FindElement(selenium.ByCSSSelector, ".normal_price")
			if err != nil {
				return steam_helper.Trace(err, costMainElement)
			}

			fmt.Println("======= 11 =========")

			costStr, err := costElement.GetAttribute("data-price")
			if err != nil {
				return steam_helper.Trace(err, costElement)
			}

			fmt.Println("======= 12 =========")

			countElement, err := skin.FindElement(selenium.ByCSSSelector, ".market_listing_num_listings_qty")
			if err != nil {
				return steam_helper.Trace(err, skin)
			}

			fmt.Println("======= 13 =========")

			countStr, err := countElement.Text()
			if err != nil {
				return steam_helper.Trace(err, countElement)
			}

			fmt.Println("======= 14 =========")

			nextBtn, err := wd.FindElement(selenium.ByCSSSelector, ".pagebtn")
			if err != nil {
				return steam_helper.Trace(err, wd)
			}

			fmt.Println("======= 15 =========")

			end, err := steam_helper.MoveMouseAndClick(nextBtn, start)
			if err != nil {
				return steam_helper.Trace(err, nextBtn)
			}
			start = end

			fmt.Println("======= 16 =========")

			cost, err := strconv.Atoi(costStr)
			if err != nil {
				return steam_helper.Trace(err)
			}

			fmt.Println("======= 17 =========")

			count, err := strconv.Atoi(countStr)
			if err != nil {
				return steam_helper.Trace(err)
			}

			fmt.Println("======= 18 =========")

			steamSkins = append(steamSkins, entity.SeleniumSteamSkin{
				HashName: hashName,
				RuName:   ruName,
				Cost:     cost,
				Count:    count,
			})
		}

		fmt.Println("======= 19 =========")
		ch.WriteModel(context.Background(), steamSkins)
	}
}
