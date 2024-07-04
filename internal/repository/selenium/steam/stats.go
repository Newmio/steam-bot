package reposteam

import (
	"bot/internal/domain/entity"
	"context"
	"fmt"
	"strconv"

	"github.com/Newmio/steam_helper"

	"github.com/tebeka/selenium"
)

func (r *steam) CheckTradeItem(wd selenium.WebDriver, links []string, ch steam_helper.CursorCh[entity.CheckItem]) {

	for _, link := range links {
		if err := wd.Get(link); err != nil {
			ch.WriteError(context.Background(), steam_helper.Trace(err))
			return
		}

		steam_helper.SleepRandom(1000, 2000)

		element, err := wd.FindElement(selenium.ByCSSSelector, ".searchResultsRows")
		if err != nil {
			ch.WriteError(context.Background(), steam_helper.Trace(err, wd))
			return
		}

		items, err := element.FindElements(selenium.ByCSSSelector, ".market_listing_row.market_recent_listing_row")
		if err != nil {
			ch.WriteError(context.Background(), steam_helper.Trace(err, wd))
			return
		}

		fmt.Println(len(items))

		steam_helper.SleepRandom(1000, 2000)
	}
}

func (r *steam) SynchItems(wd selenium.WebDriver, game string, ch steam_helper.CursorCh[[]entity.SteamItem]) {
	var url string
	page, stop := 1, 0

	switch game {
	case "csgo":
		url = "https://steamcommunity.com/market/search?appid=730#p1_popular_desc"
	}

	if err := wd.Get(url); err != nil {
		ch.WriteError(context.Background(), steam_helper.Trace(err, wd))
		return
	}

	steam_helper.SleepRandom(1000, 2000)

	start, err := steam_helper.GetStartMousePosition(wd)
	if err != nil {
		ch.WriteError(context.Background(), steam_helper.Trace(err, wd))
		return
	}

	for {
		var steamSkins []entity.SteamItem

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

			countStr, err := countElement.GetAttribute("data-qty")
			if err != nil {
				ch.WriteError(context.Background(), steam_helper.Trace(err, countElement))
				return
			}

			link, err := skin.GetAttribute("href")
			if err != nil {
				ch.WriteError(context.Background(), steam_helper.Trace(err, skin))
				return
			}

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

			steamSkins = append(steamSkins, entity.SteamItem{
				HashName: hashName,
				RuName:   ruName,
				Cost:     cost,
				Count:    count,
				Link:     link,
			})
		}

		ch.WriteModel(context.Background(), steamSkins)

		if page == 0 {
			stopPages, err := wd.FindElements(selenium.ByCSSSelector, ".market_paging_pagelink")
			if err != nil {
				ch.WriteError(context.Background(), steam_helper.Trace(err, wd))
				return
			}

			for _, value := range stopPages {
				p, err := value.Text()
				if err != nil {
					ch.WriteError(context.Background(), steam_helper.Trace(err, value))
					return
				}

				intP, err := strconv.Atoi(p)
				if err != nil {
					ch.WriteError(context.Background(), steam_helper.Trace(err))
					return
				}

				if intP > stop {
					stop = intP
				}
			}
		}

		if stop == page {
			return
		} else {
			page++
		}

		nextBtn, err := wd.FindElements(selenium.ByCSSSelector, ".pagebtn")
		if err != nil {
			ch.WriteError(context.Background(), steam_helper.Trace(err, wd))
			return
		}

		end, err := steam_helper.MoveMouseAndClick(wd, nextBtn[1], start)
		if err != nil {
			ch.WriteError(context.Background(), steam_helper.Trace(err, nextBtn))
			return
		}
		start = end

		steam_helper.SleepRandom(4000, 6000)
	}
}
