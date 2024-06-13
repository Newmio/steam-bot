package reposteam

import (
	"bot/internal/domain/entity"
	"fmt"
	"strconv"

	"github.com/Newmio/steam_helper"

	"github.com/tebeka/selenium"
)

func (r *steam) GetCSGOStats(wd selenium.WebDriver, ch chan entity.SteamSkin) error {
	page := 1

	if err := wd.Get(fmt.Sprintf("https://steamcommunity.com/market/search?appid=730#p%d_popular_desc", page)); err != nil {
		return steam_helper.Trace(err, wd)
	}

	steam_helper.SleepRandom(9000, 10000)

	start, err := steam_helper.GetRandomStartMousePosition(wd)
	if err != nil {
		return steam_helper.Trace(err, wd)
	}

	for {

		skin, err := wd.FindElement(selenium.ByCSSSelector, ".market_listing_row_link")
		if err != nil {
			return steam_helper.Trace(err, wd)
		}

		hashNameElement, err := skin.FindElement(selenium.ByCSSSelector, "market_listing_row.market_recent_listing_row.market_listing_searchresult")
		if err != nil {
			return steam_helper.Trace(err, skin)
		}

		hashName, err := hashNameElement.GetAttribute("hash-name")
		if err != nil {
			return steam_helper.Trace(err, hashNameElement)
		}

		ruNameElement, err := skin.FindElement(selenium.ByCSSSelector, "market_listing_item_name")
		if err != nil {
			return steam_helper.Trace(err, skin)
		}

		ruName, err := ruNameElement.Text()
		if err != nil {
			return steam_helper.Trace(err, ruNameElement)
		}

		costElement, err := skin.FindElement(selenium.ByCSSSelector, "normal_price")
		if err != nil {
			return steam_helper.Trace(err, skin)
		}

		costStr, err := costElement.GetAttribute("data-price")
		if err != nil {
			return steam_helper.Trace(err, costElement)
		}

		countElement, err := skin.FindElement(selenium.ByCSSSelector, "market_listing_num_listings_qty")
		if err != nil {
			return steam_helper.Trace(err, skin)
		}

		countStr, err := countElement.Text()
		if err != nil {
			return steam_helper.Trace(err, countElement)
		}

		nextBtn, err := wd.FindElement(selenium.ByCSSSelector, "searchResults_btn_next")
		if err != nil {
			return steam_helper.Trace(err, wd)
		}

		end, err := steam_helper.MoveMouseAndClick(nextBtn, start)
		if err != nil {
			return steam_helper.Trace(err, nextBtn)
		}
		start = end

		cost, err := strconv.Atoi(costStr)
		if err != nil {
			return steam_helper.Trace(err)
		}

		count, err := strconv.Atoi(countStr)
		if err != nil {
			return steam_helper.Trace(err)
		}

		ch <- entity.SteamSkin{
			HashName: hashName,
			RuName:   ruName,
			Cost:     cost,
			Count:    count,
		}

		page++
	}
}
