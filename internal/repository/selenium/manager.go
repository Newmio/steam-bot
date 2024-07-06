package reposelenium

import (
	"bot/internal/domain/entity"
	"context"
	"fmt"
	"sync"
	"time"

	repocsmoney "bot/internal/repository/selenium/csmoney"
	repodmarket "bot/internal/repository/selenium/dmarket"
	reposteam "bot/internal/repository/selenium/steam"

	"github.com/Newmio/steam_helper"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

type ISelenium interface {
	SteamLogin() error
	SynchItems(game string, ch steam_helper.CursorCh[[]entity.SteamItem])
	Ping(url string) (string, error)
	CheckTradeItems(links []string, ch steam_helper.CursorCh[entity.CheckItem])
	GetHistoryItems(links []string, ch steam_helper.CursorCh[entity.SteamSellHistory])
	GetHistoryItem(link string) (entity.SteamSellHistory, error)
}

type seleniumRepo struct {
	wd      map[string]selenium.WebDriver
	user    entity.SteamUser
	steam   reposteam.ISteam
	dmarket repodmarket.IDmarket
	csmoney repocsmoney.ICsmoney
	mu      sync.Mutex
}

func NewSelenium(user entity.SteamUser) ISelenium {
	steam_helper.BuzierOffset = 200
	steam_helper.BuzierSteps = 30

	wd, err := createDriver()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &seleniumRepo{
		csmoney: repocsmoney.NewCsmoney(),
		dmarket: repodmarket.NewDmarket(),
		steam:   reposteam.NewSteam(),
		user:    user,
		wd:      map[string]selenium.WebDriver{"steam": wd},
	}
}

func (r *seleniumRepo) GetHistoryItem(link string) (entity.SteamSellHistory, error) {
	ch := make(steam_helper.CursorCh[entity.SteamSellHistory])

	wd, err := r.getDriver("steam")
	if err != nil {
		return entity.SteamSellHistory{}, steam_helper.Trace(err)
	}

	go r.steam.GetHistoryItems(wd, []string{link}, ch)

	select {
	case history, ok := <-ch:
		if !ok {
			return entity.SteamSellHistory{}, nil
		}
		if history.Error != nil {
			return entity.SteamSellHistory{}, steam_helper.Trace(history.Error)
		}

		return history.Model, nil

	case <-time.After(time.Minute):
		return entity.SteamSellHistory{}, steam_helper.Trace(fmt.Errorf("timeout"))
	}
}

func (r *seleniumRepo) GetHistoryItems(links []string, ch steam_helper.CursorCh[entity.SteamSellHistory]) {
	wd, err := r.getDriver("steam")
	if err != nil {
		ch.WriteError(context.Background(), steam_helper.Trace(err))
		return
	}

	r.steam.GetHistoryItems(wd, links, ch)
}

func (r *seleniumRepo) CheckTradeItems(links []string, ch steam_helper.CursorCh[entity.CheckItem]) {
	wd, err := r.getDriver("steam")
	if err != nil {
		ch.WriteError(context.Background(), steam_helper.Trace(err))
		return
	}

	r.steam.CheckTradeItems(wd, links, ch)
}

func (r *seleniumRepo) SynchItems(game string, ch steam_helper.CursorCh[[]entity.SteamItem]) {
	wd, err := r.getDriver("steam")
	if err != nil {
		ch.WriteError(context.Background(), steam_helper.Trace(err))
		return
	}

	r.steam.SynchItems(wd, game, ch)
}

func (r *seleniumRepo) SteamLogin() error {
	wd, err := r.getDriver("steam")
	if err != nil {
		return steam_helper.Trace(err)
	}

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
		r.mu.Unlock()
	}

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

		//Extensions: []string{"proxyauth.zip"},
	}

	chromeCaps.AddExtension("proxy_auth_plugin.zip")
	chromeCaps.AddExtension("delete_headers_plugin.zip")
	chromeCaps.AddExtension("JJICBEFPEMNPHINCCGIKPDAAGJEBBNHG_4_3_1_0.crx")

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
