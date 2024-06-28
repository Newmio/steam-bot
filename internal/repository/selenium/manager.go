package reposelenium

import (
	"bot/internal/domain/entity"
	"fmt"
	"sync"
	"time"

	repodmarket "bot/internal/repository/selenium/dmarket"
	reposteam "bot/internal/repository/selenium/steam"

	"github.com/Newmio/steam_helper"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

type ISelenium interface {
	SteamLogin() error
	SynchSteamCSGOSkins(ch steam_helper.CursorCh[[]entity.SeleniumSteamSkin]) error
	SynchDmarketCSGOSkins(ch steam_helper.CursorCh[[]entity.SeleniumSteamSkin]) error
	Ping(url string) error
}

type seleniumRepo struct {
	wd      map[string]selenium.WebDriver
	user    entity.SteamUser
	steam   reposteam.ISteam
	dmarket repodmarket.IDmarket
	mu sync.Mutex
}

func NewSelenium(user entity.SteamUser) ISelenium {
	steam_helper.BuzierOffset = 200
	steam_helper.BuzierSteps = 30

	wd, err := createDriver()
	if err != nil {
		return nil
	}

	fmt.Println("@@@@@@@@@@@@@@@@@")
	fmt.Println(wd)
	fmt.Println("@@@@@@@@@@@@@@@@@")

	return &seleniumRepo{
		dmarket: repodmarket.NewDmarket(),
		steam:   reposteam.NewSteam(),
		user:    user,
		wd:      map[string]selenium.WebDriver{"steam": wd},
	}
}

func (r *seleniumRepo) SynchDmarketCSGOSkins(ch steam_helper.CursorCh[[]entity.SeleniumSteamSkin]) error {
	wd, err := r.getDriver("dmarket")
	if err != nil{
		return steam_helper.Trace(err)
	}

	return r.dmarket.SynchCSGOSkins(wd, ch)
}

func (r *seleniumRepo) SynchSteamCSGOSkins(ch steam_helper.CursorCh[[]entity.SeleniumSteamSkin]) error {
	wd, err := r.getDriver("steam")
	if err != nil{
		return steam_helper.Trace(err)
	}

	return r.steam.SynchCSGOSkins(wd, ch)
}

func (r *seleniumRepo) SteamLogin() error {
	wd, err := r.getDriver("steam")
	if err != nil{
		return steam_helper.Trace(err)
	}

	_, err = r.steam.Login(wd, r.user)
	if err != nil {
		return steam_helper.Trace(err)
	}

	return nil
}

func (r *seleniumRepo) getDriver(name string)(selenium.WebDriver, error){
	if _, ok := r.wd[name]; !ok{
		wd, err := createDriver()
		if err != nil{
			return nil, steam_helper.Trace(err)
		}

		r.mu.Lock()
		r.wd[name] = wd
		r.mu.Unlock()
	}

	return r.wd[name], nil
}

func (r *seleniumRepo) Ping(url string) error {
	if err := r.wd["steam"].Get(url); err != nil {
		return steam_helper.Trace(err)
	}

	time.Sleep(2 * time.Second)

	html, err := r.wd["steam"].PageSource()
	if err != nil {
		return steam_helper.Trace(err)
	}

	fmt.Println("-----------------------------------------------")
	fmt.Println("-----------------------------------------------")
	fmt.Println(html)
	fmt.Println("-----------------------------------------------")
	fmt.Println("-----------------------------------------------")

	return nil
}

func (r *seleniumRepo) Test(wd selenium.WebDriver) {

	if err := wd.Get("https://webhook.site/701a5ccc-db3d-4f79-a764-49da9674c64d"); err != nil {
		fmt.Println(steam_helper.Trace(err))
	}

	time.Sleep(time.Second * 2)

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
	agent := steam_helper.GetRandomUserAgent()
	window := steam_helper.GetRandomWindowSize()

	chromeCaps := chrome.Capabilities{
		Args: []string{
			"--disable-webgl",         // Отключение WebGL
			"--disable-webrtc",        // Отключение WebRTC
			"--disable-notifications", // Отключение уведомлений
			"--disable-rtc-smoothness-algorithm",
			"--incognito",             // Режим инкогнито
			"--lang=ru",               // Изменение языка
			"--no-sandbox",            // Отключение песочницы
			"--disable-dev-shm-usage", // Отключение использования shared memory
			"--disable-blink-features=AutomationControlled", // Отключение автоматических контролируемых функций
			"--headless",
			"--user-agent=" + agent,
			fmt.Sprintf("--window-size=%d,%d", window.Width, window.Height),
		},

		Prefs: map[string]interface{}{
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

// func createDriver(user entity.SteamUser) (selenium.WebDriver, error) {
// 	proxy := user.Proxy[rand.Intn(len(user.Proxy))]
// 	agent := steam_helper.GetRandomUserAgent()
// 	window := steam_helper.GetRandomWindowSize()

// 	//host := fmt.Sprintf("http://%s:%s@%s:%s", proxy.Login, proxy.Password, proxy.Ip, proxy.Port)

// 	firefoxPrefs := map[string]interface{}{
// 		"general.useragent.override": agent,
// 		// "webgl.disabled":                             true,  // Отключение WebGL
// 		// "media.peerconnection.enabled":               false, // Отключение WebRTC
// 		// "dom.webnotifications.enabled":               false, // Отключение уведомлений
// 		// "media.navigator.video.enabled":              false, // Отключение видео через RTC
// 		// "media.navigator.audio.enabled":              false, // Отключение аудио через RTC
// 		// "intl.accept_languages":                      "ru",  // Установка языка интерфейса
// 		// "dom.webdriver.enabled":                      false, // Отключение свойства `navigator.webdriver`
// 		// "datareporting.healthreport.uploadEnabled":   false,
// 		// "datareporting.policy.dataSubmissionEnabled": false,
// 		// "social.enabled":                             false,
// 		// "network.prefetch-next":                      false,
// 		"network.proxy.type":                      1,
// 		"network.proxy.http":                      proxy.Ip,
// 		"network.proxy.http_port":                 proxy.Port,
// 		"network.proxy.ssl":                       proxy.Ip,
// 		"network.proxy.ssl_port":                  proxy.Port,
// 		"signon.autologin.proxy":                  true,
// 		"network.proxy.allow_hijacking_localhost": true,
// 	}

// 	firefoxArgs := []string{
// 		//"--private", // Incognito mode
// 	}

// 	firefoxCaps := firefox.Capabilities{
// 		Args:  firefoxArgs,
// 		Prefs: firefoxPrefs,
// 	}

// 	caps := selenium.Capabilities{"browserName": "firefox"}
// 	caps.AddFirefox(firefoxCaps)

// 	wd, err := selenium.NewRemote(caps, "http://127.0.0.1:4444")
// 	if err != nil {
// 		return nil, steam_helper.Trace(err)
// 	}

// 	if err := wd.ResizeWindow("", window.Width, window.Height); err != nil {
// 		return nil, steam_helper.Trace(err)
// 	}

// 	script := `
//     // Функция для генерации случайного целого числа в заданном диапазоне
//     function getRandomInt(min, max) {
//         min = Math.ceil(min);
//         max = Math.floor(max);
//         return Math.floor(Math.random() * (max - min + 1)) + min;
//     }

//     // Функция для генерации случайного элемента из массива
//     function getRandomArrayElement(arr) {
//         return arr[Math.floor(Math.random() * arr.length)];
//     }

//     // Рандомизация платформы (операционной системы)
//     const platforms = ['Win32', 'Linux x86_64', 'MacIntel'];
//     Object.defineProperty(navigator, 'platform', {
//         get: function() {
//             return getRandomArrayElement(platforms);
//         }
//     });

//     // Рандомизация количества ядер процессора
//     Object.defineProperty(navigator, 'hardwareConcurrency', {
//         get: function() {
//             return getRandomInt(2, 8); // Случайное количество ядер CPU от 2 до 8
//         }
//     });

//     // Рандомизация объема памяти
//     Object.defineProperty(navigator, 'deviceMemory', {
//         get: function() {
//             return getRandomInt(4, 16); // Случайное количество памяти от 4 до 16 GB
//         }
//     });
// 	`

// 	_, err = wd.ExecuteScript(script, nil)
// 	if err != nil {
// 		return nil, steam_helper.Trace(err)
// 	}

// 	return wd, nil
// }
