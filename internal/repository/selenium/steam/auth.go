package reposteam

import (
	"bot/internal/domain/entity"
	"fmt"
	"strings"

	"github.com/Newmio/steam_helper"
	"github.com/tebeka/selenium"
)

func (r *steam) Login(wd selenium.WebDriver, user entity.SteamUser) (string, error) {
	haveCookie := true

	if err := wd.Get("https://steamcommunity.com/login/home/?goto="); err != nil {
		return "", steam_helper.Trace(err)
	}

	steam_helper.SleepRandom(10000, 15000)

	cookieMsg, err := wd.FindElement(selenium.ByCSSSelector, ".btn_blue_steamui.btn_medium.replyButton")
	if err != nil {
		haveCookie = false
	}

	inputs, err := wd.FindElements(selenium.ByCSSSelector, "._2GBWeup5cttgbTw8FM3tfx")
	if err != nil {
		fmt.Println(err)
		return "", steam_helper.Trace(err, wd)
	}

	loginBtn, err := wd.FindElement(selenium.ByCSSSelector, ".DjSvCZoKKfoNSmarsEcTS")
	if err != nil {
		return "", steam_helper.Trace(err, wd)
	}

	start, err := steam_helper.GetStartMousePosition(wd)
	if err != nil {
		return "", steam_helper.Trace(err, wd)
	}

	end, err := steam_helper.MoveMouseAndWriteText(wd, inputs[0], start, user.Login)
	if err != nil {
		return "", steam_helper.Trace(err, wd)
	}

	if haveCookie {
		end, err = steam_helper.MoveMouseAndClick(wd, cookieMsg, end)
		if err != nil {
			return "", steam_helper.Trace(err, wd)
		}
	}

	end, err = steam_helper.MoveMouseAndWriteText(wd, inputs[1], end, user.Password)
	if err != nil {
		return "", steam_helper.Trace(err, wd)
	}

	_, err = steam_helper.MoveMouseAndClick(wd, loginBtn, end)
	if err != nil {
		return "", steam_helper.Trace(err, wd)
	}

	steam_helper.SleepRandom(6000, 8000)

	link, err := wd.CurrentURL()
	if err != nil {
		return "", steam_helper.Trace(err, wd)
	}

	if !strings.Contains(link, "steamcommunity.com/profiles") {
		return "", steam_helper.Trace(fmt.Errorf("auth error"), wd)
	}

	return strings.Replace(link, "/home", "", -1), nil
}
