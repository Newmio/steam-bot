package reposteam

import (
	"bot/internal/domain/entity"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Newmio/steam_helper"

	"github.com/tebeka/selenium"
)

func (r *steam) GetHistoryItems(wd selenium.WebDriver, links []string, ch steam_helper.CursorCh[entity.SteamSellHistory]) {
	var resp map[string]interface{}

	for _, link := range links {
		if err := wd.Get(link); err != nil {
			ch.WriteError(context.Background(), steam_helper.Trace(err))
			return
		}

		pre, err := wd.FindElement(selenium.ByTagName, "pre")
		if err != nil {
			ch.WriteError(context.Background(), steam_helper.Trace(err, wd))
			return
		}

		body, err := pre.Text()
		if err != nil {
			ch.WriteError(context.Background(), steam_helper.Trace(err, pre))
			return
		}

		if err := json.Unmarshal([]byte(body), &resp); err != nil {
			ch.WriteError(context.Background(), steam_helper.Trace(err, body))
			return
		}

		var prices []entity.SteamItemPrice

		for _, value := range resp["prices"].([]interface{}){
			value2 := value.([]interface{})
			
			dateTime, err := time.Parse("Jan 02 2006 15:04", strings.Replace(value2[0].(string), " +0", "00", -1))
			if err != nil {
				ch.WriteError(context.Background(), steam_helper.Trace(err))
				return
			}

			cost := int(math.Round(value2[1].(float64) * 100))

			count, err := strconv.Atoi(fmt.Sprint(value2[2].(string)))
			if err != nil {
				ch.WriteError(context.Background(), steam_helper.Trace(err))
				return
			}

			prices = append(prices, entity.SteamItemPrice{
				DateTime: dateTime,
				Cost: cost,
				Count: count,
			})
		}

		steamSellHistory := entity.SteamSellHistory{
			PriceSuffix: resp["price_suffix"].(string),
			Prices: prices,
		}

		ch.WriteModel(context.Background(), steamSellHistory)
	}
	close(ch)
}

func (r *steam) CheckTradeItems(wd selenium.WebDriver, links []string, ch steam_helper.CursorCh[entity.CheckItem]) {

	for _, link := range links {
		if err := wd.Get(link); err != nil {
			ch.WriteError(context.Background(), steam_helper.Trace(err))
			return
		}

		steam_helper.SleepRandom(1000, 2000)

		element, err := wd.FindElement(selenium.ByCSSSelector, ".market_content_block.market_home_listing_table.market_home_main_listing_table.market_listing_table")
		if err != nil {
			ch.WriteError(context.Background(), steam_helper.Trace(err, wd))
			return
		}

		_, err = element.FindElement(selenium.ByCSSSelector, ".market_listing_table_message")
		if err == nil {
			continue
		}

		rows, err := element.FindElement(selenium.ByID, "searchResultsRows")
		if err != nil {
			ch.WriteError(context.Background(), steam_helper.Trace(err, element))
			return
		}

		items, err := rows.FindElements(selenium.ByCSSSelector, ".market_listing_row.market_recent_listing_row")
		if err != nil {
			ch.WriteError(context.Background(), steam_helper.Trace(err, rows))
			return
		}

		var costItemsSell []string
		for _, item := range items {

			costElement, err := item.FindElement(selenium.ByCSSSelector, ".market_listing_price.market_listing_price_with_fee")
			if err != nil {
				ch.WriteError(context.Background(), steam_helper.Trace(err, item))
				return
			}

			costStr, err := costElement.Text()
			if err != nil {
				ch.WriteError(context.Background(), steam_helper.Trace(err, costElement))
				return
			}

			costItemsSell = append(costItemsSell, costStr)
		}

		btnDiv, err := wd.FindElement(selenium.ByID, "market_buyorder_info_show_details")
		if err != nil {
			ch.WriteError(context.Background(), steam_helper.Trace(err, wd))
			return
		}

		btn, err := btnDiv.FindElement(selenium.ByTagName, "span")
		if err != nil {
			ch.WriteError(context.Background(), steam_helper.Trace(err, btnDiv))
			return
		}

		start, err := steam_helper.GetStartMousePosition(wd)
		if err != nil {
			ch.WriteError(context.Background(), steam_helper.Trace(err))
			return
		}

		if _, err := steam_helper.MoveMouseAndClick(wd, btn, start); err != nil {
			ch.WriteError(context.Background(), steam_helper.Trace(err, wd))
			return
		}

		table, err := wd.FindElement(selenium.ByCSSSelector, ".market_commodity_orders_table")
		if err != nil {
			ch.WriteError(context.Background(), steam_helper.Trace(err, wd))
			return
		}

		orderElements, err := table.FindElements(selenium.ByTagName, "tr")
		if err != nil {
			ch.WriteError(context.Background(), steam_helper.Trace(err, table))
			return
		}

		itemsBuy := make(map[string]string)

		for i, orderElement := range orderElements {
			if i == 0 {
				continue
			}

			order, err := orderElement.FindElements(selenium.ByTagName, "td")
			if err != nil {
				ch.WriteError(context.Background(), steam_helper.Trace(err, orderElement))
				return
			}

			costStr, err := order[0].Text()
			if err != nil {
				ch.WriteError(context.Background(), steam_helper.Trace(err, order[0]))
				return
			}

			countStr, err := order[1].Text()
			if err != nil {
				ch.WriteError(context.Background(), steam_helper.Trace(err, order[1]))
				return
			}

			itemsBuy[costStr] = countStr
		}

		link, err := wd.CurrentURL()
		if err != nil {
			ch.WriteError(context.Background(), steam_helper.Trace(err))
			return
		}

		sell := make(map[int]int)
		buy := make(map[int]int)

		re := regexp.MustCompile("[^0-9]")

		for _, cost := range costItemsSell {
			costClear := re.ReplaceAllString(cost, "")

			costInt, err := strconv.Atoi(costClear)
			if err != nil {
				ch.WriteError(context.Background(), steam_helper.Trace(err))
				return
			}

			if strings.Contains(cost, ",") {
				sell[costInt]++
			}else{
				sell[costInt*100]++
			}
		}

		for cost, count := range itemsBuy {
			costClear := re.ReplaceAllString(cost, "")

			costInt, err := strconv.Atoi(costClear)
			if err != nil {
				ch.WriteError(context.Background(), steam_helper.Trace(err))
				return
			}

			countInt, err := strconv.Atoi(count)
			if err != nil {
				ch.WriteError(context.Background(), steam_helper.Trace(err))
				return
			}

			if strings.Contains(cost, ",") {
				buy[costInt] = countInt
			}else{
				buy[costInt*100] = countInt
			}
		}

		hashName := strings.Split(link, "/")

		checkItem := entity.CheckItem{
			HashName: hashName[len(hashName)-1],
			Sell: sell,
			Buy:  buy,
		}

		ch.WriteModel(context.Background(), checkItem)

		steam_helper.SleepRandom(1000, 2000)
	}
	close(ch)
}

func (r *steam) SynchItems(wd selenium.WebDriver, game string, ch steam_helper.CursorCh[[]entity.SteamItem]) {
	var reqUrl string
	page, stop := 1, 0

	switch game {
	case "csgo":
		reqUrl = "https://steamcommunity.com/market/search?appid=730#p1_popular_desc"
	}

	if err := wd.Get(reqUrl); err != nil {
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

			link, err := skin.GetAttribute("href")
			if err != nil {
				ch.WriteError(context.Background(), steam_helper.Trace(err, skin))
				return
			}

			imgLinkElement, err := hashNameElement.FindElement(selenium.ByTagName, "img")
			if err != nil {
				ch.WriteError(context.Background(), steam_helper.Trace(err, hashNameElement))
				return
			}

			imgLink, err := imgLinkElement.GetAttribute("src")
			if err != nil {
				ch.WriteError(context.Background(), steam_helper.Trace(err, imgLinkElement))
				return
			}

			steamSkins = append(steamSkins, entity.SteamItem{
				HashName: url.QueryEscape(hashName),
				RuName:   ruName,
				Link:     link,
				ImgLink:  imgLink+"2x",
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
			close(ch)
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

		steam_helper.SleepRandom(3000, 5000)
	}
}
