package reposelenium

import (
	"bot/internal/domain/entity"
	"fmt"

	reposteam "bot/internal/repository/selenium/steam"

	"github.com/Newmio/steam_helper"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

type ISelenium interface {
	SteamLogin(user entity.SteamUser) error
	GetCSGOSkins(login string, ch steam_helper.CursorCh[[]entity.SeleniumSteamSkin]) error
}

type seleniumRepo struct {
	wd    selenium.WebDriver
	user  entity.SteamUser
	steam reposteam.ISteam
}

func NewSelenium(user entity.SteamUser) ISelenium {
	steam_helper.BuzierOffset = 200
	steam_helper.BuzierSteps = 30

	//proxy := user.Proxy[rand.Intn(len(user.Proxy))]
	//proxyAddress := fmt.Sprintf("http://%s:%s@%s:%s", proxy.Login, proxy.Password, proxy.Ip, proxy.Port)
	//proxyAddress := "http://yggdjocl:vajq3n53awr1@38.154.227.167:5868"
	//agent := steam_helper.GetRandomUserAgent()
	//window := steam_helper.GetRandomWindowSize(agent)

	// extensionData, err := io.ReadFile("proxyauth.zip")
	// if err != nil {
	// 	log.Fatalf("Ошибка чтения файла расширения: %v", err)
	// }

	chromeCaps := chrome.Capabilities{
		Args: []string{
			"--disable-webgl",         // Отключение WebGL
			"--disable-webrtc",        // Отключение WebRTC
			"--disable-notifications", // Отключение уведомлений
			"--incognito",             // Режим инкогнито
			// "--lang=en-US",            // Изменение языка
			"--no-sandbox",            // Отключение песочницы
			"--disable-dev-shm-usage", // Отключение использования shared memory
			"--disable-blink-features=AutomationControlled", // Отключение автоматических контролируемых функций
			// "--headless",
			//"--user-agent=" + agent,
			//fmt.Sprintf("--window-size=%d,%d", window.Width, window.Height),
			"--window-size=1920,1080",
			"--user-agent=TEST",
			//"--proxy-server=" + proxyAddress,
		},

		Prefs: map[string]interface{}{
			"proxy": map[string]interface{}{
				"mode":          "fixed_servers",
				"server":        "38.154.227.167:5868",
				"proxyType":     "MANUAL",
				"httpProxy":     "38.154.227.167:5868",
				"sslProxy":      "38.154.227.167:5868",
				"socksUsername": "yggdjocl",
				"socksPassword": "vajq3n53awr1",
				"noProxy":       "",
				"autodetect":    false,
				"class":         "org.openqa.selenium.Proxy",
			},

			"profile.default_content_setting_values.notifications": 0,
			"profile.default_content_setting_values.images":        0,
			"profile.managed_default_content_settings.popups":      0,
		},

		//Extensions: []string{"proxyauth.zip"},
	}

	caps := selenium.Capabilities{"browserName": "chrome"}
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

func (r *seleniumRepo) GetCSGOSkins(login string, ch steam_helper.CursorCh[[]entity.SeleniumSteamSkin]) error {
	return r.steam.GetCSGOSkins(r.wd, ch)
}

func (r *seleniumRepo) SteamLogin(user entity.SteamUser) error {
	// r.Test(r.wd)

	// return nil

	_, err := r.steam.Login(r.wd, user)
	if err != nil {
		return steam_helper.Trace(err)
	}

	return nil
}

func (r *seleniumRepo) Test(wd selenium.WebDriver) {

	if err := wd.Get("https://webhook.site/84b33734-6ca5-4133-b18a-7fda14bb1f04"); err != nil {
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
