package reposelenium

import (
	"bot/internal/domain/entity"
	"context"
	"fmt"
	"sync"
	"time"

	repocsmoney "bot/internal/repository/selenium/csmoney"
	repodmarket "bot/internal/repository/selenium/dmarket"
	repohelpers "bot/internal/repository/selenium/helpers"
	reposteam "bot/internal/repository/selenium/steam"

	"github.com/Newmio/steam_helper"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

type ISelenium interface {
	SteamLogin() error
	SynchItems(info entity.PaginationInfo[[]entity.SteamItem])
	Ping(url string) (string, error)
	CheckItems(info entity.PaginationInfo[entity.CheckItem])
	GetHistoryItems(info entity.PaginationInfo[[]entity.SteamSellHistory])
	GetHistoryItem(link string) ([]entity.SteamSellHistory, error)
	GetLinksForTradeItem(game string) (map[string]float64, error)
	GetRareFloats(limit, offset int) (map[string][]entity.FloatItem, error)
}

type seleniumRepo struct {
	wd      map[string]selenium.WebDriver
	wg      map[string]*sync.WaitGroup
	user    entity.SteamUser
	steam   reposteam.ISteam
	dmarket repodmarket.IDmarket
	csmoney repocsmoney.ICsmoney
	helpers repohelpers.IHelpers
	mu      sync.Mutex
}

func NewSelenium(user entity.SteamUser) ISelenium {
	steam_helper.BuzierOffset = 200
	steam_helper.BuzierSteps = 30

	return &seleniumRepo{
		helpers: repohelpers.NewHelpers(),
		csmoney: repocsmoney.NewCsmoney(user.Proxy),
		dmarket: repodmarket.NewDmarket(),
		steam:   reposteam.NewSteam(),
		user:    user,
		wd:      make(map[string]selenium.WebDriver),
		wg:      make(map[string]*sync.WaitGroup),
	}
}

func (r *seleniumRepo) GetRareFloats(limit, offset int) (map[string][]entity.FloatItem, error) {
	return r.csmoney.GetRareFloats(limit, offset)
}

func (r *seleniumRepo) GetLinksForTradeItem(game string) (map[string]float64, error) {
	wd, err := r.getDriver("steam")
	if err != nil {
		return nil, steam_helper.Trace(err)
	}
	defer r.wg["steam"].Done()

	return r.helpers.GetLinksForTradeItem(wd, game)
}

func (r *seleniumRepo) GetHistoryItem(link string) ([]entity.SteamSellHistory, error) {
	info := entity.PaginationInfo[[]entity.SteamSellHistory]{
		Links: []string{link},
		Ch:    make(steam_helper.CursorCh[[]entity.SteamSellHistory]),
	}

	wd, err := r.getDriver("helpers")
	if err != nil {
		return nil, steam_helper.Trace(err)
	}
	defer r.wg["helpers"].Done()

	go r.steam.GetHistoryItems(wd, info)

	select {
	case history, ok := <-info.Ch:
		if !ok {
			return nil, nil
		}
		if history.Error != nil {
			return nil, steam_helper.Trace(history.Error)
		}

		return history.Model, nil

	case <-time.After(time.Minute):
		return nil, steam_helper.Trace(fmt.Errorf("timeout"))
	}
}

func (r *seleniumRepo) GetHistoryItems(info entity.PaginationInfo[[]entity.SteamSellHistory]) {
	wd, err := r.getDriver("steam")
	if err != nil {
		info.Ch.WriteError(context.Background(), steam_helper.Trace(err))
		return
	}
	defer r.wg["steam"].Done()

	r.steam.GetHistoryItems(wd, info)
}

func (r *seleniumRepo) CheckItems(info entity.PaginationInfo[entity.CheckItem]) {
	wd, err := r.getDriver("steam")
	if err != nil {
		info.Ch.WriteError(context.Background(), steam_helper.Trace(err))
		return
	}
	defer r.wg["steam"].Done()

	r.steam.CheckItems(wd, info)
}

func (r *seleniumRepo) SynchItems(info entity.PaginationInfo[[]entity.SteamItem]) {
	wd, err := r.getDriver("steam")
	if err != nil {
		info.Ch.WriteError(context.Background(), steam_helper.Trace(err))
		return
	}
	defer r.wg["steam"].Done()

	r.steam.SynchItems(wd, info)
}

func (r *seleniumRepo) SteamLogin() error {
	wd, err := r.getDriver("steam")
	if err != nil {
		return steam_helper.Trace(err)
	}
	defer r.wg["steam"].Done()

	_, err = r.steam.Login(wd, r.user)
	if err != nil {
		return steam_helper.Trace(err)
	}

	return nil
}

func (r *seleniumRepo) getDriver(name string) (selenium.WebDriver, error) {
	if _, ok := r.wd[name]; !ok {
		wd, err := createDriver()
		if err != nil {
			return nil, steam_helper.Trace(err)
		}

		r.mu.Lock()
		r.wd[name] = wd
		r.wg[name] = new(sync.WaitGroup)
		r.mu.Unlock()

		if _, err := r.steam.Login(wd, r.user); err != nil {
			return nil, steam_helper.Trace(err)
		}
	}

	r.wg[name].Wait()
	r.wg[name].Add(1)
	return r.wd[name], nil
}

func (r *seleniumRepo) Ping(url string) (string, error) {
	wd, err := r.getDriver("steam")
	if err != nil {
		return "", steam_helper.Trace(err)
	}

	// r.Test(wd)
	// return "", nil

	if err := wd.Get(url); err != nil {
		return "", steam_helper.Trace(err)
	}

	time.Sleep(2 * time.Second)

	html, err := wd.PageSource()
	if err != nil {
		return "", steam_helper.Trace(err)
	}

	return html, nil
}

func (r *seleniumRepo) Test(wd selenium.WebDriver) {

	if err := wd.Get("https://steamcommunity.com/market/pricehistory/?appid=730&market_hash_name=AK-47%20%7C%20Redline%20%28Field-Tested%29"); err != nil {
		fmt.Println(steam_helper.Trace(err))
		return
	}

	// if _, err := wd.ExecuteScript("window.open('about:blank', '_blank');", nil); err != nil {
	// 	fmt.Println(steam_helper.Trace(err))
	// 	return
	// }

	tabs, err := wd.WindowHandles()
	if err != nil {
		fmt.Println(steam_helper.Trace(err))
		return
	}

	fmt.Println("tabs2", tabs)

	time.Sleep(2 * time.Minute)
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
			//"--incognito",             // Режим инкогнито
			"--lang=ru",               // Изменение языка
			"--no-sandbox",            // Отключение песочницы
			"--disable-dev-shm-usage", // Отключение использования shared memory
			"--disable-blink-features=AutomationControlled", // Отключение автоматических контролируемых функций
			//"--headless", // Не отображать окно браузера
			"--user-agent=" + agent,
			fmt.Sprintf("--window-size=%d,%d", window.Width, window.Height),
		},

		Prefs: map[string]interface{}{
			"profile.default_content_setting_values.notifications": 0,
			"profile.default_content_setting_values.images":        0,
			"profile.managed_default_content_settings.popups":      0,
		},
	}

	chromeCaps.AddExtension("proxy_auth_plugin.zip")
	chromeCaps.AddExtension("delete_headers_plugin.zip")
	//chromeCaps.AddExtension("JJICBEFPEMNPHINCCGIKPDAAGJEBBNHG_4_3_1_0.crx")
	//chromeCaps.AddExtension("Steam-Inventory-Helper-Chrome.zip")

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
