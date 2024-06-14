package app

import (
	"bot/internal/configs/account"
	"bot/internal/configs/proxy"
	"bot/internal/domain/entity"
	"bot/internal/domain/usecase"
	usecasesteam "bot/internal/domain/usecase/steam"
	reposelenium "bot/internal/repository/selenium"

	"bot/internal/transport/http"

	"github.com/labstack/echo/v4"
)

func Init() {
	e := echo.New()

	accounts := account.InitAccounts()
	proxy := proxy.InitProxy()

	bots := make(map[string]entity.Bot)
	for _, account := range accounts {
		bots[account.Login] = entity.Bot{SteamUser: account, Proxy: proxy}
	}

	seleniumRepo := reposelenium.NewSelenium()
	steamUsecase := usecasesteam.NewSteam(seleniumRepo)
	usecase := usecase.NewUseCase(steamUsecase, bots)
	authHandler := http.NewHandler(usecase)
	authHandler.InitRoutes(e)

	e.Logger.Fatal(e.Start(":8088"))
}
