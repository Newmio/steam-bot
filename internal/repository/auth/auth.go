package repoauth

import (
	"bot/internal/domain/entity"
	"fmt"
	"math/rand"
	"strings"

	"github.com/Newmio/steam_helper"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

type IAuth interface {
	Login(login string) (entity.AuthInfo, error)
}

type auth struct {
	wb       map[string]selenium.WebDriver
	accounts map[string][]entity.ProxyAccount
}

func NewAuth(accounts map[string][]entity.ProxyAccount) IAuth {
	return &auth{accounts: accounts, wb: make(map[string]selenium.WebDriver)}
}

func (r *auth) GetChromeDriver(login string) (selenium.WebDriver, error) {

	if r.wb[login] == nil {

		//proxy := r.accounts[login][rand.Intn(len(r.accounts[login]))]
		//proxyAddress := fmt.Sprintf("http://%s:%s@%s:%s", proxy.Login, proxy.Password, proxy.Ip, proxy.Port)
		//agent := steam_helper.GetRandomUserAgent()
		//window := steam_helper.GetRandomWindowSize(agent)

		chromeCaps := chrome.Capabilities{
			Args: []string{
				"--headless",
				"--no-sandbox",
				"--disable-dev-shm-usage",
				//"--user-agent=" + agent,
				//fmt.Sprintf("--window-size=%d,%d", window.Width, window.Height),
				"--window-size=1920,1080",
			},

			// Prefs: map[string]interface{}{
			// 	"proxy": map[string]interface{}{
			// 		"httpProxy": proxyAddress,
			// 		"sslProxy":  proxyAddress,
			// 		"proxyType": "MANUAL",
			// 	},
			// },
		}

		caps := selenium.Capabilities{"browserName": "chrome"}
		caps.AddChrome(chromeCaps)

		wb, err := selenium.NewRemote(caps, "http://127.0.0.1:9515")
		if err != nil {
			return nil, steam_helper.Trace(err)
		}

		r.wb[login] = wb
		return wb, nil
	}

	return r.wb[login], nil
}

func (r *auth) test(login string) {
	wb, err := r.GetChromeDriver(login)
	if err != nil {
		fmt.Println(steam_helper.Trace(err))
	}

	if err := wb.Get("https://ru.wikipedia.org/wiki/Заглавная_страница"); err != nil {
		fmt.Println(steam_helper.Trace(err))
	}

	btns, err := wb.FindElements(selenium.ByCSSSelector, ".main-wikimedia-listItem")
	if err != nil {
		fmt.Println(steam_helper.Trace(err))
	}

	// if err := steam_helper.MoveMouse(wb, 0, 0, btnLocation.X, btnLocation.Y); err != nil {
	// 	fmt.Println(steam_helper.Trace(err))
	// }

	steam_helper.SleepRandom(1000, 2000)

	_, err = steam_helper.MoveMouseAndClick(btns[0], steam_helper.Position{X: 0, Y: 0})
	if err != nil {
		fmt.Println(steam_helper.Trace(err))
	}

	steam_helper.SleepRandom(1000, 2000)

	if err := btns[0].Click(); err != nil {
		fmt.Println(steam_helper.Trace(err))
	}
}

func (r *auth) Login(login string) (entity.AuthInfo, error) {

	// r.test(login)

	// return entity.AuthInfo{}, nil

	var authInfo entity.AuthInfo
	accounts := r.accounts[login]
	acc := accounts[rand.Intn(len(accounts))]

	wb, err := r.GetChromeDriver(login)
	if err != nil {
		return authInfo, steam_helper.Trace(err)
	}

	if err := wb.Get("https://steamcommunity.com/login/home/?goto="); err != nil {
		return authInfo, steam_helper.Trace(err)
	}

	steam_helper.SleepRandom(4000, 5000)

	cookieMsg, err := wb.FindElement(selenium.ByCSSSelector, ".btn_blue_steamui.btn_medium.replyButton")
	if err != nil {
		if !strings.Contains(err.Error(), "no such element") {
			return authInfo, steam_helper.Trace(err)
		}
	}

	inputs, err := wb.FindElements(selenium.ByCSSSelector, "._2eKVn6g5Yysx9JmutQe7WV")
	if err != nil {
		fmt.Println(err)
		return authInfo, steam_helper.Trace(err)
	}

	loginBtn, err := wb.FindElement(selenium.ByCSSSelector, "._2QgFEj17t677s3x299PNJQ")
	if err != nil {
		return authInfo, steam_helper.Trace(err)
	}

	start, err := steam_helper.GetRandomStartMousePosition(wb)
	if err != nil {
		return authInfo, steam_helper.Trace(err)
	}

	end, err := steam_helper.MoveMouseAndClick(cookieMsg, start)
	if err != nil {
		return authInfo, steam_helper.Trace(err)
	}

	end, err = steam_helper.MoveMouseAndWriteText(inputs[0], end, acc.SteamLogin)
	if err != nil {
		return authInfo, steam_helper.Trace(err)
	}

	end, err = steam_helper.MoveMouseAndWriteText(inputs[1], end, acc.SteamPassword)
	if err != nil {
		return authInfo, steam_helper.Trace(err)
	}

	_, err = steam_helper.MoveMouseAndClick(loginBtn, end)
	if err != nil {
		return authInfo, steam_helper.Trace(err)
	}

	steam_helper.SleepRandom(4000, 5000)

	profile, err := wb.FindElement(selenium.ByCSSSelector, ".user_avatar.playerAvatar.offline")
	if err != nil {
		return authInfo, steam_helper.Trace(err)
	}

	link, err := profile.GetAttribute("href")
	if err != nil {
		return authInfo, steam_helper.Trace(err)
	}

	fmt.Println(link)

	if !strings.Contains(link, "steamcommunity.com/profiles") {
		return authInfo, fmt.Errorf("auth error")
	}

	return authInfo, nil
}

// func (r *auth) Login2(login string) (entity.AuthInfo, error) {

// 	r.test(login)

// 	return entity.AuthInfo{}, nil

// 	var authInfo entity.AuthInfo
// 	accounts := r.accounts[login]
// 	acc := accounts[rand.Intn(len(accounts))]

// 	wb, err := r.GetChromeDriver(login)
// 	if err != nil {
// 		return authInfo, steam_helper.Trace(err)
// 	}
// 	defer wb.Quit()

// 	if err := wb.Get("https://steamcommunity.com/login/home/?goto="); err != nil {
// 		return authInfo, steam_helper.Trace(err)
// 	}

// 	steam_helper.SleepRandom(4000, 5000)

// 	// cookieMsg, err := wb.FindElement(selenium.ByCSSSelector, ".btn_blue_steamui.btn_medium.replyButton")
// 	// if err != nil {
// 	// 	if !strings.Contains(err.Error(), "no such element") {
// 	// 		return authInfo, steam_helper.Trace(err)
// 	// 	}
// 	// }

// 	inputs, err := wb.FindElements(selenium.ByCSSSelector, "._2eKVn6g5Yysx9JmutQe7WV")
// 	if err != nil {
// 		fmt.Println(err)
// 		return authInfo, steam_helper.Trace(err)
// 	}

// 	loginBtn, err := wb.FindElement(selenium.ByCSSSelector, "._2QgFEj17t677s3x299PNJQ")
// 	if err != nil {
// 		return authInfo, steam_helper.Trace(err)
// 	}

// 	loginLocation, err := inputs[0].Location()
// 	if err != nil {
// 		return authInfo, steam_helper.Trace(err)
// 	}

// 	passLocation, err := inputs[1].Location()
// 	if err != nil {
// 		return authInfo, steam_helper.Trace(err)
// 	}

// 	loginBtnLocation, err := loginBtn.Location()
// 	if err != nil {
// 		return authInfo, steam_helper.Trace(err)
// 	}

// 	window, err := wb.ExecuteScript("return [window.innerWidth, window.innerHeight];", nil)
// 	if err != nil {
// 		return authInfo, steam_helper.Trace(err)
// 	}
// 	windowSize := window.([]interface{})

// 	loginSize, err := inputs[0].Size()
// 	if err != nil {
// 		return authInfo, steam_helper.Trace(err)
// 	}

// 	passSize, err := inputs[1].Size()
// 	if err != nil {
// 		return authInfo, steam_helper.Trace(err)
// 	}

// 	loginBtnSize, err := loginBtn.Size()
// 	if err != nil {
// 		return authInfo, steam_helper.Trace(err)
// 	}

// 	startX := int(windowSize[0].(float64))
// 	startY := int(windowSize[1].(float64))

// 	loginEndX := loginLocation.X + rand.Intn(loginSize.Width)
// 	passEndX := passLocation.X + rand.Intn(passSize.Width)
// 	loginBtnEndX := loginBtnLocation.X + rand.Intn(loginBtnSize.Width)

// 	loginEndY := loginLocation.Y + rand.Intn(loginSize.Height)
// 	passEndY := passLocation.Y + rand.Intn(passSize.Height)
// 	loginBtnEndY := loginBtnLocation.Y + rand.Intn(loginBtnSize.Height)

// 	if err := steam_helper.MoveMouse(wb, startX, startY, loginEndX, loginEndY); err != nil {
// 		return authInfo, steam_helper.Trace(err)
// 	}

// 	steam_helper.SleepRandom(100, 500)

// 	if err := inputs[0].Click(); err != nil {
// 		return authInfo, steam_helper.Trace(err)
// 	}

// 	steam_helper.SleepRandom(100, 500)

// 	for _, value := range acc.SteamLogin {
// 		if err := inputs[0].SendKeys(string(value)); err != nil {
// 			return authInfo, steam_helper.Trace(err)
// 		}
// 		steam_helper.SleepRandom(50, 150)
// 	}

// 	if err := steam_helper.MoveMouse(wb, loginEndX, loginEndY, passEndX, passEndY); err != nil {
// 		return authInfo, steam_helper.Trace(err)
// 	}

// 	steam_helper.SleepRandom(100, 500)

// 	if err := inputs[1].Click(); err != nil {
// 		return authInfo, steam_helper.Trace(err)
// 	}

// 	steam_helper.SleepRandom(100, 500)

// 	for _, value := range acc.SteamPassword {
// 		if err := inputs[1].SendKeys(string(value)); err != nil {
// 			return authInfo, steam_helper.Trace(err)
// 		}
// 		steam_helper.SleepRandom(50, 150)
// 	}

// 	steam_helper.SleepRandom(100, 500)

// 	if err := steam_helper.MoveMouse(wb, passEndX, passEndY, loginBtnEndX, loginBtnEndY); err != nil {
// 		return authInfo, steam_helper.Trace(err)
// 	}

// 	steam_helper.SleepRandom(100, 500)

// 	if err := loginBtn.Click(); err != nil {
// 		return authInfo, steam_helper.Trace(err)
// 	}

// 	steam_helper.SleepRandom(4000, 5000)

// 	profile, err := wb.FindElement(selenium.ByCSSSelector, ".user_avatar.playerAvatar.offline")
// 	if err != nil {
// 		return authInfo, steam_helper.Trace(err)
// 	}

// 	link, err := profile.GetAttribute("href")
// 	if err != nil {
// 		return authInfo, steam_helper.Trace(err)
// 	}

// 	fmt.Println(link)

// 	if !strings.Contains(link, "steamcommunity.com/profiles") {
// 		return authInfo, fmt.Errorf("auth error")
// 	}

// 	return authInfo, nil
// }
