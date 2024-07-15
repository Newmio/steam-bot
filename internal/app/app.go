package app

import (
	"archive/zip"
	"bot/internal/configs/bot"
	"bot/internal/configs/redis"
	"bot/internal/configs/sqlite"
	"bot/internal/domain/usecase"
	usecasedmarket "bot/internal/domain/usecase/dmarket"
	usecasesteam "bot/internal/domain/usecase/steam"
	repodb "bot/internal/repository/db"
	reporedis "bot/internal/repository/db/redis"
	reposqlite "bot/internal/repository/db/sqlite"
	reposelenium "bot/internal/repository/selenium"
	"fmt"
	"io"
	"os"

	"bot/internal/transport/http"

	"github.com/labstack/echo/v4"
)

func Init() {
	e := echo.New()
	botConfig := bot.Init()

	sqlite, err := sqlite.OpenDb(botConfig.Bot.SteamUser.Login)
	if err != nil {
		panic(err)
	}

	redis, err := redis.OpenDb()
	if err != nil {
		panic(err)
	}

	if err := buildProxyExtension(botConfig.Bot.SteamUser.Proxy[0].Ip, botConfig.Bot.SteamUser.Proxy[0].Port,
		botConfig.Bot.SteamUser.Proxy[0].Login, botConfig.Bot.SteamUser.Proxy[0].Password); err != nil {
		panic(err)
	}

	if err := buildDeleteHeadersExtension(); err != nil {
		panic(err)
	}

	seleniumRepo := reposelenium.NewSelenium(botConfig.Bot.SteamUser)
	repoRedis := reporedis.NewRedis(redis)
	repoSqlite := reposqlite.NewSqlite(sqlite)
	dbRepo := repodb.NewDatabase(repoRedis, repoSqlite)
	steamUsecase := usecasesteam.NewSteam(seleniumRepo, dbRepo, botConfig.Markets["steam"])
	dmarketUsecase := usecasedmarket.NewDmarket(seleniumRepo, dbRepo)
	usecase := usecase.NewUseCase(steamUsecase, dmarketUsecase, botConfig.Bot)
	authHandler := http.NewHandler(usecase)
	authHandler.InitRoutes(e)

	e.Logger.Fatal(e.Start(":8088"))
}

func buildDeleteHeadersExtension() error {

	manifest_json := `{
  "manifest_version": 2,
  "name": "Remove Headers",
  "version": "1.0",
  "permissions": ["webRequest", "webRequestBlocking", "<all_urls>"],
  "background": {
    "scripts": ["background.js"]
  }
}`

	background_js := `chrome.webRequest.onBeforeSendHeaders.addListener(
  function(details) {
    const headersToRemove = ["sec-ch-ua-platform"];
    details.requestHeaders = details.requestHeaders.filter(header => !headersToRemove.includes(header.name));
    return { requestHeaders: details.requestHeaders };
  },
  { urls: ["<all_urls>"] },
  ["blocking", "requestHeaders"]
);`

	return createZip(manifest_json, background_js, "delete_headers_plugin.zip")
}

func buildProxyExtension(host, port, userName, password string) error {

	manifest_json := `{
		  "version": "1.0.0",
		  "manifest_version": 2,
		  "name": "Chrome Proxy",
		  "permissions": [
			"proxy",
			"tabs",
			"unlimitedStorage",
			"storage",
			"<all_urls>",
			"webRequest",
			"webRequestBlocking"
		  ],
		  "background": {
			"scripts": ["background.js"]
		  },
		  "minimum_chrome_version":"22.0.0"
		}`

	background_js := fmt.Sprintf(`var config = {
		  mode: "fixed_servers",
		  rules: {
			singleProxy: {
			  scheme: "http",
			  host: "%s",
			  port: parseInt(%s)
			},
			bypassList: ["localhost"]
		  }
		};
		
		chrome.proxy.settings.set({value: config, scope: "regular"}, function() {});
		
		function callbackFn(details) {
		  return {
			authCredentials: {
			  username: "%s",
			  password: "%s"
			}
		  };
		}
		
		chrome.webRequest.onAuthRequired.addListener(
		  callbackFn,
		  {urls: ["<all_urls>"]},
		  ['blocking']
		);`, host, port, userName, password)

	return createZip(manifest_json, background_js, "proxy_auth_plugin.zip")
}

func createZip(manifest_json, background_js, fileName string) error {

	fos, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("os.Create Error: %w", err)
	}
	defer fos.Close()
	zipWriter := zip.NewWriter(fos)

	preManifestFile, err := os.Create("manifest.json")
	if err != nil {
		return fmt.Errorf("os.Create manifestFile Error: %w", err)
	}
	if _, err = preManifestFile.Write([]byte(manifest_json)); err != nil {
		return fmt.Errorf("preManifestFile.Write Error: %w", err)
	}
	preManifestFile.Close()
	manifestFile, err := os.Open("manifest.json")
	if err != nil {
		return fmt.Errorf("os.Open manifestFile  Error: %w", err)
	}

	wmf, err := zipWriter.Create("manifest.json")
	if err != nil {
		return fmt.Errorf("zipWriter.Create Error: %w", err)
	}
	if _, err := io.Copy(wmf, manifestFile); err != nil {
		return fmt.Errorf("io.Copy Error: %w", err)
	}
	manifestFile.Close()

	preBackFile, err := os.Create("background.js")
	if err != nil {
		return fmt.Errorf("os.Create preBackFile  Error: %w", err)
	}
	if _, err = preBackFile.Write([]byte(background_js)); err != nil {
		return fmt.Errorf("preBackFile.Write Error: %w", err)
	}
	preBackFile.Close()

	backFile, err := os.Open("background.js")
	if err != nil {
		return fmt.Errorf("os.Open backFile Error: %w", err)
	}

	wbf, err := zipWriter.Create("background.js")
	if err != nil {
		return fmt.Errorf("zipWriter.Create(background.js) Error: %w", err)
	}
	if _, err := io.Copy(wbf, backFile); err != nil {
		return fmt.Errorf("io.Copy Error: %w", err)
	}
	backFile.Close()
	if err := os.Remove("manifest.json"); err != nil {
		return fmt.Errorf("os.Remove(manifest.json) Error: %w", err)
	}

	if err := os.Remove("background.js"); err != nil {
		return fmt.Errorf("os.Remove(background.js) Error: %w", err)
	}

	return zipWriter.Close()
}
