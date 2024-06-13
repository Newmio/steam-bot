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
	SteamLogin(login string) error
}

type seleniumRepo struct {
	wd       map[string]selenium.WebDriver
	accounts map[string][]entity.ProxyAccount
	auth     reposteam.IAuth
}

func NewSelenium(accounts map[string][]entity.ProxyAccount) ISelenium {
	return &seleniumRepo{accounts: accounts, auth: reposteam.NewAuth()}
}

func (r *seleniumRepo) SteamLogin(login string) error {
	wd, err := r.getChromeDriver(login)
	if err != nil {
		return steam_helper.Trace(err)
	}

	return r.auth.Login(wd, r.getSteamUser(login))
}

func (r *seleniumRepo) getSteamUser(login string) entity.User {
	acc := r.accounts[login][0]

	return entity.User{Login: acc.SteamLogin, Password: acc.SteamPassword}
}

func (r *seleniumRepo) getChromeDriver(login string) (selenium.WebDriver, error) {
	if len(r.accounts[login]) == 0 {
		return nil, steam_helper.Trace(fmt.Errorf("no accounts for login %s", login))
	}

	if r.wd[login] == nil {

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
