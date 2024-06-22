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
	SynchSteamCSGOSkins(login string, ch steam_helper.CursorCh[[]entity.SeleniumSteamSkin]) error
}

type seleniumRepo struct {
	wd    map[string]selenium.WebDriver
	user  entity.SteamUser
	steam reposteam.ISteam
}

func NewSelenium(user entity.SteamUser) ISelenium {
	steam_helper.BuzierOffset = 200
	steam_helper.BuzierSteps = 30

	wd, err := createDriver()
	if err != nil {
		panic(err)
	}

	return &seleniumRepo{
		steam: reposteam.NewSteam(),
		user:  user,
		wd:    map[string]selenium.WebDriver{"steam": wd},
	}
}

func (r *seleniumRepo) SynchSteamCSGOSkins(login string, ch steam_helper.CursorCh[[]entity.SeleniumSteamSkin]) error {
	return r.steam.SynchCSGOSkins(r.wd["steam"], ch)
}

func (r *seleniumRepo) SteamLogin(user entity.SteamUser) error {
	// r.Test(r.wd)

	// return nil

	_, err := r.steam.Login(r.wd["steam"], user)
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

func createDriver() (selenium.WebDriver, error) {
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
			"--disable-rtc-smoothness-algorithm",
			"--incognito",             // Режим инкогнито
			// "--lang=en-US",            // Изменение языка
			"--no-sandbox",            // Отключение песочницы
			"--disable-dev-shm-usage", // Отключение использования shared memory
			"--disable-blink-features=AutomationControlled", // Отключение автоматических контролируемых функций
			"--headless",
			//"--user-agent=" + agent,
			//fmt.Sprintf("--window-size=%d,%d", window.Width, window.Height),
			"--window-size=1920,1080",
			//"--user-agent=TEST",
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

	script := `
    // Функция для генерации случайного целого числа в заданном диапазоне
    function getRandomInt(min, max) {
        min = Math.ceil(min);
        max = Math.floor(max);
        return Math.floor(Math.random() * (max - min + 1)) + min;
    }

    // Функция для генерации случайного элемента из массива
    function getRandomArrayElement(arr) {
        return arr[Math.floor(Math.random() * arr.length)];
    }

    // Рандомизация платформы (операционной системы)
    const platforms = ['Win32', 'Linux x86_64', 'MacIntel'];
    Object.defineProperty(navigator, 'platform', {
        get: function() {
            return getRandomArrayElement(platforms);
        }
    });

    // Рандомизация количества ядер процессора
    Object.defineProperty(navigator, 'hardwareConcurrency', {
        get: function() {
            return getRandomInt(2, 8); // Случайное количество ядер CPU от 2 до 8
        }
    });

    // Рандомизация объема памяти
    Object.defineProperty(navigator, 'deviceMemory', {
        get: function() {
            return getRandomInt(4, 16); // Случайное количество памяти от 4 до 16 GB
        }
    });
	`

	caps := selenium.Capabilities{"browserName": "chrome"}
	caps.AddChrome(chromeCaps)

	wd, err := selenium.NewRemote(caps, "http://127.0.0.1:9515")
	if err != nil {
		return nil, steam_helper.Trace(err)
	}

	_, err = wd.ExecuteScript(script, nil)
    if err != nil {
        return nil, steam_helper.Trace(err)
    }

	return wd, nil
}
