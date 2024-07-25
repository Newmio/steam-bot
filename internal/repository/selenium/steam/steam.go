package reposteam

import (
	"bot/internal/domain/entity"
	"strings"

	"github.com/Newmio/steam_helper"
	"github.com/tebeka/selenium"
)

type ISteam interface {
	Login(wd selenium.WebDriver, user entity.SteamUser) (string, error)
	SynchItems(wd selenium.WebDriver, info entity.PaginationInfo[[]entity.SteamItem])
	CheckItems(wd selenium.WebDriver, info entity.PaginationInfo[entity.CheckItem])
	GetHistoryItems(wd selenium.WebDriver, info entity.PaginationInfo[[]entity.SteamSellHistory])
}

type steam struct{}

func NewSteam() ISteam {
	return &steam{}
}

func (r *steam) ifTooManyRequests(wd selenium.WebDriver) (bool, error) {
	var msg string

	link, err := wd.CurrentURL()
	if err != nil {
		return false, steam_helper.Trace(err)
	}

	if strings.Contains(link, "steamcommunity.com/market/pricehistory") {

		msgElement, err := wd.FindElement(selenium.ByTagName, "pre")
		if err != nil {
			return false, nil
		}

		msg, err = msgElement.Text()
		if err != nil {
			return false, steam_helper.Trace(err)
		}

	} else {
		msgElement, err := wd.FindElement(selenium.ByID, "message")
		if err != nil {
			return false, nil
		}

		h3, err := msgElement.FindElement(selenium.ByTagName, "h3")
		if err != nil {
			return false, steam_helper.Trace(err)
		}

		msg, err = h3.Text()
		if err != nil {
			return false, steam_helper.Trace(err)
		}
	}

	if msg != "Вы делали слишком много запросов. Пожалуйста, подождите и повторите запрос позже." && msg != "null"{
		return false, nil
	}

	if err := wd.Get("https://steamcommunity.com"); err != nil {
		return false, steam_helper.Trace(err)
	}

	steam_helper.SleepRandom(1000, 2000)

	homeTabs, err := wd.FindElement(selenium.ByCSSSelector, ".community_home_tabs")
	if err != nil {
		return false, steam_helper.Trace(err)
	}

	a, err := homeTabs.FindElements(selenium.ByTagName, "a")
	if err != nil {
		return false, steam_helper.Trace(err)
	}

	start, err := steam_helper.GetStartMousePosition(wd)
	if err != nil {
		return false, steam_helper.Trace(err)
	}

	end, err := steam_helper.MoveMouseAndClick(wd, a[6], start)
	if err != nil {
		return false, steam_helper.Trace(err)
	}

	steam_helper.SleepRandom(500, 700)

	navigation, err := wd.FindElement(selenium.ByCSSSelector, ".supernav_container")
	if err != nil {
		return false, steam_helper.Trace(err)
	}

	if _, err := steam_helper.MoveMouseAndClick(wd, navigation, end); err != nil {
		return false, steam_helper.Trace(err)
	}

	steam_helper.SleepRandom(1000, 2000)

	return true, nil
}
