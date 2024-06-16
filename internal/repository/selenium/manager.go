package reposelenium

import (
	"bot/internal/domain/entity"
	"fmt"
	"math/rand"

	reposteam "bot/internal/repository/selenium/steam"

	"github.com/Newmio/steam_helper"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

type ISelenium interface {
	SteamLogin(user entity.SteamUser) error
	GetCSGOSkins(login string, ch chan []entity.SteamSkin) error
}

type seleniumRepo struct {
	wd    selenium.WebDriver
	user  entity.SteamUser
	steam reposteam.ISteam
}

func NewSelenium(user entity.SteamUser) ISelenium {
	steam_helper.BuzierOffset = 200
	steam_helper.BuzierSteps = 30

	proxy := user.Proxy[rand.Intn(len(user.Proxy))]
	proxyAddress := fmt.Sprintf("http://%s:%s@%s:%s", proxy.Login, proxy.Password, proxy.Ip, proxy.Port)
	//agent := steam_helper.GetRandomUserAgent()
	//window := steam_helper.GetRandomWindowSize(agent)

	fmt.Println(proxyAddress)

	chromeCaps := chrome.Capabilities{
		Args: []string{
			// "--disable-webgl",         // Отключение WebGL
			// "--disable-webrtc",        // Отключение WebRTC
			// "--disable-notifications", // Отключение уведомлений
			// "--incognito",             // Режим инкогнито
			// "--lang=en-US",            // Изменение языка
			// "--no-sandbox",            // Отключение песочницы
			// "--disable-dev-shm-usage", // Отключение использования shared memory
			// "--disable-blink-features=AutomationControlled", // Отключение автоматических контролируемых функций
			// "--headless",
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

	// p := selenium.Proxy{
	// 	Type: selenium.Manual,
	// 	HTTP: proxyAddress,
	// 	SSL:  proxyAddress,
	// }

	caps := selenium.Capabilities{"browserName": "chrome"} //"proxy": p}
	caps.AddChrome(chromeCaps)

	wd, err := selenium.NewRemote(caps, "http://127.0.0.1:9515")
	if err != nil {
		panic(err)
	}

	return &seleniumRepo{
		steam: reposteam.NewSteam(),
		user:  user,
		wd:    wd,
	}
}

func (r *seleniumRepo) GetCSGOSkins(login string, ch chan []entity.SteamSkin) error {
	return r.steam.GetCSGOSkins(r.wd, ch)
}

func (r *seleniumRepo) SteamLogin(user entity.SteamUser) error {
	// r.Test(wd)

	// return nil

	_, err := r.steam.Login(r.wd, user)
	if err != nil {
		return steam_helper.Trace(err)
	}

	return nil
}

// func (r *seleniumRepo) getChromeDriver(login string) (selenium.WebDriver, error) {

// 	if wd, ok := r.wd[login]; !ok || wd == nil {

// 		proxy := r.users[login].Proxy[rand.Intn(len(r.users[login].Proxy))]
// 		proxyAddress := fmt.Sprintf("http://%s:%s@%s:%s", proxy.Login, proxy.Password, proxy.Ip, proxy.Port)
// 		//agent := steam_helper.GetRandomUserAgent()
// 		//window := steam_helper.GetRandomWindowSize(agent)

// 		fmt.Println(proxyAddress)

// 		chromeCaps := chrome.Capabilities{
// 			Args: []string{
// 				// "--headless",
// 				// "--no-sandbox",
// 				// "--disable-dev-shm-usage",
// 				//"--user-agent=" + agent,
// 				//fmt.Sprintf("--window-size=%d,%d", window.Width, window.Height),
// 				"--window-size=1920,1080",
// 			},

// 			// Prefs: map[string]interface{}{
// 			// 	"proxy": map[string]interface{}{
// 			// 		"httpProxy": proxyAddress,
// 			// 		"sslProxy":  proxyAddress,
// 			// 		"proxyType": "MANUAL",
// 			// 	},
// 			// },
// 		}

// 		// p := selenium.Proxy{
// 		// 	Type: selenium.Manual,
// 		// 	HTTP: proxyAddress,
// 		// 	SSL:  proxyAddress,
// 		// }

// 		caps := selenium.Capabilities{"browserName": "chrome"} //"proxy": p}
// 		caps.AddChrome(chromeCaps)

// 		wd, err := selenium.NewRemote(caps, "http://127.0.0.1:9515")
// 		if err != nil {
// 			return nil, steam_helper.Trace(err)
// 		}

// 		r.wd[login] = wd
// 		return wd, nil
// 	}

// 	return r.wd[login], nil
// }

func (r *seleniumRepo) Test(wd selenium.WebDriver) {

	if err := wd.Get("https://webhook.site/a85394f2-323e-4faa-834c-6ad90cb62754"); err != nil {
		fmt.Println(steam_helper.Trace(err))
	}

	// btns, err := wd.FindElements(selenium.ByCSSSelector, ".left-eye")
	// if err != nil {
	// 	fmt.Println(steam_helper.Trace(err))
	// }

	// steam_helper.SleepRandom(1000, 2000)

	// // start, err := steam_helper.GetStartMousePosition(wd)
	// // if err != nil {
	// // 	fmt.Println(steam_helper.Trace(err))
	// // }

	// _, err = steam_helper.TestMoveMouseAndClick(wd, btns[0], steam_helper.Position{X: 100, Y: 100})
	// if err != nil {
	// 	fmt.Println(steam_helper.Trace(err))
	// }

	// steam_helper.SleepRandom(1000, 2000)
}
