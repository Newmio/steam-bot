package repohelpers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Newmio/steam_helper"
	"github.com/tebeka/selenium"
)

func (h *helpers) GetLinksForTradeItem(wd selenium.WebDriver, game string) (map[string]float64, error) {
	var url string
	var rangeClick int
	links := make(map[string]float64)

	switch game {
	case "csgo":
		url = "https://skins-table.xyz/table/"
		rangeClick = 17
	case "dota2":
		url = "https://skins-table.xyz/table_dota/"
		rangeClick = 7
	default:
		return nil, steam_helper.Trace(fmt.Errorf("unknown game"))
	}

	if err := wd.Get(url); err != nil {
		return nil, steam_helper.Trace(err)
	}

	steam_helper.SleepRandom(3000, 5000)

	steamLogin, err := wd.CurrentURL()
	if err != nil {
		return nil, steam_helper.Trace(err)
	}

	start, err := steam_helper.GetStartMousePosition(wd)
	if err != nil {
		return nil, steam_helper.Trace(err)
	}

	if strings.Contains(steamLogin, "steamcommunity.com") {
		btn, err := wd.FindElement(selenium.ByCSSSelector, ".btn_green_white_innerfade")
		if err != nil {
			return nil, steam_helper.Trace(err, wd)
		}
		end, err := steam_helper.MoveMouseAndClick(wd, btn, start)
		if err != nil {
			return nil, steam_helper.Trace(err, btn)
		}
		start = end

		steam_helper.SleepRandom(7000, 9000)

		if err := wd.Get(url); err != nil {
			return nil, steam_helper.Trace(err)
		}

		steam_helper.SleepRandom(7000, 9000)
	}

	rightBtn, err := wd.FindElement(selenium.ByID, "scroll-right")
	if err != nil {
		return nil, steam_helper.Trace(err, wd)
	}

	table, err := wd.FindElement(selenium.ByCSSSelector, ".table.table-bordered")
	if err != nil {
		return nil, steam_helper.Trace(err, wd)
	}

	offRefresh, err := table.FindElement(selenium.ByCSSSelector, ".tgl-btn")
	if err != nil {
		return nil, steam_helper.Trace(err, table)
	}

	start, err = steam_helper.MoveMouseAndClick(wd, offRefresh, start)
	if err != nil {
		return nil, steam_helper.Trace(err, offRefresh)
	}

	for range rangeClick {
		end, err := steam_helper.MoveMouseAndClick(wd, rightBtn, start)
		if err != nil {
			return nil, steam_helper.Trace(err, rightBtn)
		}
		start = end
		steam_helper.SleepRandom(1000, 1100)
	}

	sitesFirst, err := wd.FindElement(selenium.ByCSSSelector, ".sites.first")
	if err != nil {
		return nil, steam_helper.Trace(err, wd)
	}

	divsFirst, err := sitesFirst.FindElements(selenium.ByTagName, "div")
	if err != nil {
		return nil, steam_helper.Trace(err, sitesFirst)
	}

	for _, value := range divsFirst {
		marketName, err := value.GetAttribute("data-name")
		if err != nil {
			return nil, steam_helper.Trace(err, value)
		}

		if marketName == "SteamCommunity(AUTO)" {
			end, err := steam_helper.MoveMouseAndClick(wd, value, start)
			if err != nil {
				return nil, steam_helper.Trace(err, value)
			}
			start = end
			break
		}
	}

	sitestSecond, err := wd.FindElement(selenium.ByCSSSelector, ".sites.second")
	if err != nil {
		return nil, steam_helper.Trace(err, wd)
	}

	divsSecond, err := sitestSecond.FindElements(selenium.ByTagName, "div")
	if err != nil {
		return nil, steam_helper.Trace(err, sitestSecond)
	}

	for _, value := range divsSecond {
		marketName, err := value.GetAttribute("data-name")
		if err != nil {
			return nil, steam_helper.Trace(err, value)
		}

		if marketName == "SteamCommunity" {
			end, err := steam_helper.MoveMouseAndClick(wd, value, start)
			if err != nil {
				return nil, steam_helper.Trace(err, value)
			}
			start = end
			break
		}
	}

	n2, err := table.FindElement(selenium.ByID, "n2")
	if err != nil {
		return nil, steam_helper.Trace(err, table)
	}

	n2.Clear()

	start, err = steam_helper.MoveMouseAndWriteText(wd, n2, start, "300")
	if err != nil {
		return nil, steam_helper.Trace(err, n2)
	}

	per1, err := table.FindElement(selenium.ByID, "per1_to")
	if err != nil {
		return nil, steam_helper.Trace(err, table)
	}

	start, err = steam_helper.MoveMouseAndWriteText(wd, per1, start, "60")
	if err != nil {
		return nil, steam_helper.Trace(err, per1)
	}

	sort, err := table.FindElement(selenium.ByCSSSelector, ".sort")
	if err != nil {
		return nil, steam_helper.Trace(err, table)
	}

	sortTh, err := sort.FindElements(selenium.ByTagName, "th")
	if err != nil {
		return nil, steam_helper.Trace(err, sort)
	}

	for _, value := range sortTh {
		sortMode, err := value.GetAttribute("onclick")
		if err != nil {
			return nil, steam_helper.Trace(err, value)
		}

		if sortMode == "setsort(4);" {
			end, err := steam_helper.MoveMouseAndClick(wd, value, start)
			if err != nil {
				return nil, steam_helper.Trace(err, value)
			}
			start = end
			break
		}
	}

	thead, err := table.FindElement(selenium.ByTagName, "thead")
	if err != nil {
		return nil, steam_helper.Trace(err, table)
	}

	theadTr, err := thead.FindElements(selenium.ByTagName, "tr")
	if err != nil {
		return nil, steam_helper.Trace(err, thead)
	}

	var theadTh []selenium.WebElement

	for _, value := range theadTr {
		class, err := value.GetAttribute("class")
		if err != nil {
			return nil, steam_helper.Trace(err, value)
		}

		if class != "sort" {
			th, err := value.FindElements(selenium.ByTagName, "th")
			if err != nil {
				return nil, steam_helper.Trace(err, value)
			}
			theadTh = append(theadTh, th...)
		}
	}

	for _, value := range theadTh {
		colspan, err := value.GetAttribute("colspan")
		if err != nil {
			if !strings.Contains(err.Error(), "nil return value") {
				return nil, steam_helper.Trace(err, value)
			}
		}

		if colspan == "3" {
			btn, err := value.FindElement(selenium.ByCSSSelector, ".reset-button")
			if err != nil {
				return nil, steam_helper.Trace(err, value)
			}

			end, err := steam_helper.MoveMouseAndClick(wd, btn, start)
			if err != nil {
				return nil, steam_helper.Trace(err, value)
			}
			start = end
			break
		}
	}

	if game == "csgo" {
		time.Sleep(time.Second)

		window, err := wd.FindElement(selenium.ByCSSSelector, ".modal-content")
		if err != nil {
			return nil, steam_helper.Trace(err, wd)
		}

		inputForm, err := window.FindElement(selenium.ByCSSSelector, ".row-fluid")
		if err != nil {
			return nil, steam_helper.Trace(err, window)
		}

		inputs, err := inputForm.FindElements(selenium.ByTagName, "input")
		if err != nil {
			return nil, steam_helper.Trace(err, inputForm)
		}

		remove := []string{"Souvenir", "Sticker", "Graffiti", "StatTrak"}

		for _, value := range remove {
			end, err := steam_helper.MoveMouseAndWriteText(wd, inputs[0], start, value)
			if err != nil {
				return nil, steam_helper.Trace(err, inputs[0])
			}
			start = end

			end, err = steam_helper.MoveMouseAndClick(wd, inputs[1], start)
			if err != nil {
				return nil, steam_helper.Trace(err, inputs[1])
			}
			start = end
		}

		closeButtons, err := window.FindElement(selenium.ByCSSSelector, ".close")
		if err != nil {
			return nil, steam_helper.Trace(err, inputForm)
		}

		spans, err := closeButtons.FindElements(selenium.ByTagName, "span")
		if err != nil {
			return nil, steam_helper.Trace(err, closeButtons)
		}

		for _, value := range spans {
			atr, err := value.GetAttribute("aria-hidden")
			if err != nil {
				return nil, steam_helper.Trace(err, value)
			}

			if atr == "true" {
				end, err := steam_helper.MoveMouseAndClick(wd, value, start)
				if err != nil {
					return nil, steam_helper.Trace(err, value)
				}
				start = end
				break
			}
		}
		steam_helper.SleepRandom(500, 1000)
	}

	respTable, err := table.FindElement(selenium.ByTagName, "tbody")
	if err != nil {
		return nil, steam_helper.Trace(err, table)
	}

	rows, err := respTable.FindElements(selenium.ByTagName, "tr")
	if err != nil {
		return nil, steam_helper.Trace(err, respTable)
	}

	lastTd, err := rows[len(rows)-1].FindElements(selenium.ByTagName, "td")
	if err != nil {
		return nil, steam_helper.Trace(err, rows[len(rows)-1])
	}

	_, err = steam_helper.MoveMouseAndClick(wd, lastTd[0], start)
	if err != nil {
		return nil, steam_helper.Trace(err, rows[len(rows)-1])
	}

	rows, err = respTable.FindElements(selenium.ByTagName, "tr")
	if err != nil {
		return nil, steam_helper.Trace(err, respTable)
	}

	for _, item := range rows {
		td, err := item.FindElements(selenium.ByTagName, "td")
		if err != nil {
			return nil, steam_helper.Trace(err, item)
		}

		for _, value := range td {
			aElement, err := value.FindElements(selenium.ByTagName, "a")
			if err != nil {
				return nil, steam_helper.Trace(err, value)
			}

			for _, value := range aElement {
				link, err := value.GetAttribute("href")
				if err != nil {
					return nil, steam_helper.Trace(err, value)
				}

				if strings.Contains(link, "steamcommunity.com") {
					profitSpan, err := td[3].FindElement(selenium.ByTagName, "span")
					if err != nil {
						return nil, steam_helper.Trace(err, td[3])
					}

					profitStr, err := profitSpan.Text()
					if err != nil {
						return nil, steam_helper.Trace(err, profitSpan)
					}

					profit, err := strconv.ParseFloat(strings.Replace(profitStr, "%", "", -1), 64)
					if err != nil {
						return nil, steam_helper.Trace(err, profitSpan)
					}

					links[link] = profit
				}
			}
		}
	}

	return links, nil
}
