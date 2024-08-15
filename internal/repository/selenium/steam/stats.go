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
	"github.com/PuerkitoBio/goquery"

	"github.com/tebeka/selenium"
)

func (r *steam) GetHistoryItems(wd selenium.WebDriver, info entity.PaginationInfo[[]entity.SteamSellHistory]) {
	var resp map[string]interface{}
	var restartCount int

	for _, link := range info.Links {
	start:
		if err := wd.Get(link); err != nil {
			info.Ch.WriteError(context.Background(), steam_helper.Trace(err))
			return
		}

		time.Sleep(2 * time.Second)

		pre, err := wd.FindElement(selenium.ByTagName, "pre")
		if err != nil {
			info.Ch.WriteError(context.Background(), steam_helper.Trace(err, wd))
			return
		}

		body, err := pre.Text()
		if err != nil {
			info.Ch.WriteError(context.Background(), steam_helper.Trace(err, pre))
			return
		}

		if err := json.Unmarshal([]byte(body), &resp); err != nil {
			info.Ch.WriteError(context.Background(), steam_helper.Trace(err, body))
			return
		}

		if len(resp) == 0 {
			if _, err := r.ifTooManyRequests(wd); err != nil {
				info.Ch.WriteError(context.Background(), steam_helper.Trace(err))
				return
			}
			restartCount++
			if restartCount < 5 {
				goto start
			} else {
				continue
			}
		}

		switch resp["prices"].(type) {
		case bool:
			time.Sleep(time.Second)
			restartCount++
			if restartCount < 5 {
				goto start
			} else {
				continue
			}
		}

		var history []entity.SteamSellHistory

		for _, value := range resp["prices"].([]interface{}) {
			value2 := value.([]interface{})

			dateTime, err := time.Parse("Jan 02 2006 15:04", strings.Replace(value2[0].(string), " +0", "00", -1))
			if err != nil {
				info.Ch.WriteError(context.Background(), steam_helper.Trace(err))
				return
			}

			cost := int(math.Round(value2[1].(float64) * 100))

			count, err := strconv.Atoi(fmt.Sprint(value2[2].(string)))
			if err != nil {
				info.Ch.WriteError(context.Background(), steam_helper.Trace(err))
				return
			}

			history = append(history, entity.SteamSellHistory{
				PriceSuffix: resp["price_suffix"].(string),
				Price: entity.SteamItemPrice{
					DateTime: dateTime,
					Cost:     cost,
					Count:    count,
				},
			})
		}

		info.Ch.WriteModel(context.Background(), history)
		restartCount = 0
	}
	close(info.Ch)
}

func (r *steam) CheckItems(wd selenium.WebDriver, info entity.PaginationInfo[entity.CheckItem]) {
	var checkedBadLinks bool
	var badLinks []string

start:
	for _, link := range info.Links {

		if err := wd.Get(link); err != nil {
			info.Ch.WriteError(context.Background(), steam_helper.Trace(err))
			return
		}

		steam_helper.SleepRandom(1000, 2000)

		tables, err := wd.FindElements(selenium.ByCSSSelector, ".market_commodity_orders_table_container")
		if err != nil {
			info.Ch.WriteError(context.Background(), steam_helper.Trace(err, wd))
			return
		}

		if len(tables) == 0 {
			flag, err := r.ifTooManyRequests(wd)
			if err != nil {
				info.Ch.WriteError(context.Background(), steam_helper.Trace(err))
				return
			}

			if !flag {
				badLinks = append(badLinks, link)
			}
			continue
		}

		var linksToGame []string
		var sell, buy map[int]int

		if len(tables) == 2 {
			sell, err = r.findInItemTable(tables[0])
			if err != nil {
				info.Ch.WriteError(context.Background(), steam_helper.Trace(err))
				return
			}

			buy, err = r.findInItemTable(tables[1])
			if err != nil {
				info.Ch.WriteError(context.Background(), steam_helper.Trace(err))
				return
			}

		} else {
			var checkFloat bool

			if info.CommonInfo != nil {
				checkFloat = info.CommonInfo.(bool)
			}

			buy, err = r.findInItemTable(tables[0])
			if err != nil {
				info.Ch.WriteError(context.Background(), steam_helper.Trace(err))
				return
			}

			sell, linksToGame, err = r.ifNotCommodity(wd, checkFloat)
			if err != nil {
				info.Ch.WriteError(context.Background(), steam_helper.Trace(err))
				return
			}
		}

		link, err := wd.CurrentURL()
		if err != nil {
			info.Ch.WriteError(context.Background(), steam_helper.Trace(err))
			return
		}

		hashName := strings.Split(link, "/")

		checkItem := entity.CheckItem{
			HashName: hashName[len(hashName)-1],
			LinksToGame: linksToGame,
			Sell:     sell,
			Buy:      buy,
		}

		info.Ch.WriteModel(context.Background(), checkItem)

		steam_helper.SleepRandom(1000, 2000)
	}

	if len(badLinks) > 0 && !checkedBadLinks {
		checkedBadLinks = true
		info.Links = badLinks
		goto start
	}
	close(info.Ch)
}

func (r *steam) findInItemTable(table selenium.WebElement) (map[int]int, error) {
	itemsBuy := make(map[string]string)

restartSearch:

	html, err := table.GetAttribute("outerHTML")
	if err != nil {
		if !strings.Contains(err.Error(), steam_helper.ERROR_NO_SUCH_ELEMENT_IN_FRAME) {
			return nil, steam_helper.Trace(err, table)
		}
		goto restartSearch
	}
	
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, steam_helper.Trace(err, table)
	}

	doc.Find("tr").Each(func(i int, s1 *goquery.Selection) {
		if i != 0 {
			var costStr, countStr string

			s1.Find("td").Each(func(j int, s2 *goquery.Selection) {
				if j == 0 {
					costStr = s2.Text()
				} else {
					countStr = s2.Text()
				}
			})

			itemsBuy[costStr] = countStr
		}
	})

	re := regexp.MustCompile("[^0-9]")
	buy := make(map[int]int)

	for cost, count := range itemsBuy {
		costClear := re.ReplaceAllString(cost, "")

		costInt, err := strconv.Atoi(costClear)
		if err != nil {
			return nil, steam_helper.Trace(err)
		}

		countInt, err := strconv.Atoi(count)
		if err != nil {
			return nil, steam_helper.Trace(err)
		}

		if strings.Contains(cost, ",") || strings.Contains(cost, ".") {
			buy[costInt] = countInt
		} else {
			buy[costInt*100] = countInt
		}
	}

	return buy, nil
}

func (r *steam) ifNotCommodity(wd selenium.WebDriver, checkFloat bool) (map[int]int, []string, error) {
	var checkFloatCount int
	var linksToGame []string
	sell := make(map[int]int)

	start, err := steam_helper.GetStartMousePosition(wd)
	if err != nil {
		return nil, nil, steam_helper.Trace(err)
	}

restart:

	if checkFloatCount > 10 {
		return sell, linksToGame, nil
	}

	element, err := wd.FindElement(selenium.ByCSSSelector, ".market_content_block.market_home_listing_table.market_home_main_listing_table.market_listing_table")
	if err != nil {
		return nil, nil, steam_helper.Trace(err, wd)
	}

	start, err = steam_helper.MoveMouseAndClick(wd, element, start)
	if err != nil {
		return nil, nil, steam_helper.Trace(err)
	}

	rows, err := element.FindElement(selenium.ByID, "searchResultsRows")
	if err != nil {
		return nil, nil, steam_helper.Trace(err, element)
	}

	items, err := rows.FindElements(selenium.ByCSSSelector, ".market_listing_row.market_recent_listing_row")
	if err != nil {
		return nil, nil, steam_helper.Trace(err, rows)
	}

	type itInfo struct {
		costStr string
	}

	itemInfo := make(map[int]itInfo)

	for i, item := range items {

		costElement, err := item.FindElement(selenium.ByCSSSelector, ".market_listing_price.market_listing_price_with_fee")
		if err != nil {
			return nil, nil, steam_helper.Trace(err, item)
		}

		costStr, err := costElement.Text()
		if err != nil {
			return nil, nil, steam_helper.Trace(err, costElement)
		}

		itemInfo[i] = itInfo{
			costStr: costStr,
		}
	}

	re := regexp.MustCompile("[^0-9]")

	for _, cost := range itemInfo {
		var c int

		costClear := re.ReplaceAllString(cost.costStr, "")

		if costClear != "" {
			costInt, err := strconv.Atoi(costClear)
			if err != nil {
				return nil, nil, steam_helper.Trace(err)
			}

			if strings.Contains(cost.costStr, ",") || strings.Contains(cost.costStr, ".") {
				c = costInt
				sell[c]++
			} else {
				c = costInt * 100
				sell[c]++
			}
		}
	}

	if checkFloat {
		next, err := wd.FindElement(selenium.ByID, "searchResults_btn_next")
		if err != nil {
			return nil, nil, steam_helper.Trace(err, wd)
		}

		start, err = steam_helper.MoveMouseAndClick(wd, next, start)
		if err != nil {
			return nil, nil, steam_helper.Trace(err, next)
		}

		checkFloatCount++
		steam_helper.SleepRandom(400, 500)
		goto restart
	}

	return sell, linksToGame, nil
}

func (r *steam) SynchItems(wd selenium.WebDriver, info entity.PaginationInfo[[]entity.SteamItem]) {
	var reqUrl string
	page, stop := 1, 0

	switch info.Game {
	case "csgo":
		reqUrl = fmt.Sprintf("https://steamcommunity.com/market/search?appid=730#p%d_popular_desc", info.Start)
	}

	if err := wd.Get(reqUrl); err != nil {
		info.Ch.WriteError(context.Background(), steam_helper.Trace(err, wd))
		return
	}

	steam_helper.SleepRandom(1000, 2000)

	start, err := steam_helper.GetStartMousePosition(wd)
	if err != nil {
		info.Ch.WriteError(context.Background(), steam_helper.Trace(err, wd))
		return
	}

	for {
		var steamSkins []entity.SteamItem

		skins, err := wd.FindElements(selenium.ByCSSSelector, ".market_listing_row_link")
		if err != nil {
			info.Ch.WriteError(context.Background(), steam_helper.Trace(err, wd))
			return
		}

		for _, skin := range skins {

			hashNameElement, err := skin.FindElement(selenium.ByCSSSelector, ".market_listing_row.market_recent_listing_row.market_listing_searchresult")
			if err != nil {
				info.Ch.WriteError(context.Background(), steam_helper.Trace(err, skin))
				return
			}

			hashName, err := hashNameElement.GetAttribute("data-hash-name")
			if err != nil {
				info.Ch.WriteError(context.Background(), steam_helper.Trace(err, hashNameElement))
				return
			}

			ruNameElement, err := skin.FindElement(selenium.ByCSSSelector, ".market_listing_item_name")
			if err != nil {
				info.Ch.WriteError(context.Background(), steam_helper.Trace(err, skin))
				return
			}

			ruName, err := ruNameElement.Text()
			if err != nil {
				info.Ch.WriteError(context.Background(), steam_helper.Trace(err, ruNameElement))
				return
			}

			link, err := skin.GetAttribute("href")
			if err != nil {
				info.Ch.WriteError(context.Background(), steam_helper.Trace(err, skin))
				return
			}

			imgLinkElement, err := hashNameElement.FindElement(selenium.ByTagName, "img")
			if err != nil {
				info.Ch.WriteError(context.Background(), steam_helper.Trace(err, hashNameElement))
				return
			}

			imgLink, err := imgLinkElement.GetAttribute("src")
			if err != nil {
				info.Ch.WriteError(context.Background(), steam_helper.Trace(err, imgLinkElement))
				return
			}

			steamSkins = append(steamSkins, entity.SteamItem{
				HashName: url.QueryEscape(hashName),
				RuName:   ruName,
				Link:     link,
				ImgLink:  imgLink + "2x",
			})
		}

		info.Ch.WriteModel(context.Background(), steamSkins)

		if page == 0 {
			stopPages, err := wd.FindElements(selenium.ByCSSSelector, ".market_paging_pagelink")
			if err != nil {
				info.Ch.WriteError(context.Background(), steam_helper.Trace(err, wd))
				return
			}

			for _, value := range stopPages {
				p, err := value.Text()
				if err != nil {
					info.Ch.WriteError(context.Background(), steam_helper.Trace(err, value))
					return
				}

				intP, err := strconv.Atoi(p)
				if err != nil {
					info.Ch.WriteError(context.Background(), steam_helper.Trace(err))
					return
				}

				if intP > stop {
					stop = intP
				}
			}
		}

		if page == stop || page == info.Stop {
			close(info.Ch)
			return
		} else {
			page++
		}

		nextBtn, err := wd.FindElements(selenium.ByCSSSelector, ".pagebtn")
		if err != nil {
			info.Ch.WriteError(context.Background(), steam_helper.Trace(err, wd))
			return
		}

		end, err := steam_helper.MoveMouseAndClick(wd, nextBtn[1], start)
		if err != nil {
			info.Ch.WriteError(context.Background(), steam_helper.Trace(err, nextBtn))
			return
		}
		start = end

		steam_helper.SleepRandom(3000, 5000)
	}
}
