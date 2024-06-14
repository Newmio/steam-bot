package reposelenium

import (
	"bot/internal/domain/entity"

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
	wd    map[string]selenium.WebDriver
	steam reposteam.ISteam
}

func NewSelenium() ISelenium {
	return &seleniumRepo{steam: reposteam.NewSteam()}
}

func (r *seleniumRepo) GetCSGOSkins(login string, ch chan []entity.SteamSkin) error {
	wd, err := r.getChromeDriver(login)
	if err != nil {
		return steam_helper.Trace(err)
	}

	return r.steam.GetCSGOSkins(wd, ch)
}

func (r *seleniumRepo) SteamLogin(user entity.SteamUser) error {
	wd, err := r.getChromeDriver(user.Login)
	if err != nil {
		return steam_helper.Trace(err)
	}

	_, err = r.steam.Login(wd, user)
	if err != nil {
		return steam_helper.Trace(err)
	}

	return nil
}

func (r *seleniumRepo) getChromeDriver(login string) (selenium.WebDriver, error) {

	if wd, ok := r.wd[login]; !ok || wd == nil {

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

		wd, err := selenium.NewRemote(caps, "http://127.0.0.1:9515")
		if err != nil {
			return nil, steam_helper.Trace(err)
		}

		r.wd[login] = wd
		return wd, nil
	}

	return r.wd[login], nil
}
