package reposteam

import (
	"bot/internal/domain/entity"
	"fmt"
	"strings"

	"github.com/Newmio/steam_helper"
	"github.com/tebeka/selenium"
)

type IAuth interface {
	Login(wd selenium.WebDriver, user entity.User) error
}

type auth struct{}

func NewAuth() IAuth {
	return &auth{}
}

// func (r *auth) test(login string) {
// 	wb, err := r.GetChromeDriver(login)
// 	if err != nil {
// 		fmt.Println(steam_helper.Trace(err))
// 	}

// 	if err := wb.Get("https://ru.wikipedia.org/wiki/Заглавная_страница"); err != nil {
// 		fmt.Println(steam_helper.Trace(err))
// 	}

// 	btns, err := wb.FindElements(selenium.ByCSSSelector, ".main-wikimedia-listItem")
// 	if err != nil {
// 		fmt.Println(steam_helper.Trace(err))
// 	}

// 	// if err := steam_helper.MoveMouse(wb, 0, 0, btnLocation.X, btnLocation.Y); err != nil {
// 	// 	fmt.Println(steam_helper.Trace(err))
// 	// }

// 	steam_helper.SleepRandom(1000, 2000)

// 	_, err = steam_helper.MoveMouseAndClick(btns[0], steam_helper.Position{X: 0, Y: 0})
// 	if err != nil {
// 		fmt.Println(steam_helper.Trace(err))
// 	}

// 	steam_helper.SleepRandom(1000, 2000)

// 	if err := btns[0].Click(); err != nil {
// 		fmt.Println(steam_helper.Trace(err))
// 	}
// }

func (r *auth) Login(wd selenium.WebDriver, user entity.User) error {

	if err := wd.Get("https://steamcommunity.com/login/home/?goto="); err != nil {
		return steam_helper.Trace(err)
	}

	steam_helper.SleepRandom(4000, 5000)

	cookieMsg, err := wd.FindElement(selenium.ByCSSSelector, ".btn_blue_steamui.btn_medium.replyButton")
	if err != nil {
		if !strings.Contains(err.Error(), "no such element") {
			return steam_helper.Trace(err)
		}
	}

	inputs, err := wd.FindElements(selenium.ByCSSSelector, "._2eKVn6g5Yysx9JmutQe7WV")
	if err != nil {
		fmt.Println(err)
		return steam_helper.Trace(err)
	}

	loginBtn, err := wd.FindElement(selenium.ByCSSSelector, "._2QgFEj17t677s3x299PNJQ")
	if err != nil {
		return steam_helper.Trace(err)
	}

	start, err := steam_helper.GetRandomStartMousePosition(wd)
	if err != nil {
		return steam_helper.Trace(err)
	}

	end, err := steam_helper.MoveMouseAndClick(cookieMsg, start)
	if err != nil {
		return steam_helper.Trace(err)
	}

	end, err = steam_helper.MoveMouseAndWriteText(inputs[0], end, user.Login)
	if err != nil {
		return steam_helper.Trace(err)
	}

	end, err = steam_helper.MoveMouseAndWriteText(inputs[1], end, user.Password)
	if err != nil {
		return steam_helper.Trace(err)
	}

	_, err = steam_helper.MoveMouseAndClick(loginBtn, end)
	if err != nil {
		return steam_helper.Trace(err)
	}

	steam_helper.SleepRandom(4000, 5000)

	profile, err := wd.FindElement(selenium.ByCSSSelector, ".user_avatar.playerAvatar.offline")
	if err != nil {
		return steam_helper.Trace(err)
	}

	link, err := profile.GetAttribute("href")
	if err != nil {
		return steam_helper.Trace(err)
	}

	if !strings.Contains(link, "steamcommunity.com/profiles") {
		return fmt.Errorf("auth error")
	}

	return nil
}
