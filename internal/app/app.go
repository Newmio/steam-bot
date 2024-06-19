package app

import (
	"bot/internal/configs/bot"
	"bot/internal/configs/sqlite"
	"bot/internal/domain/entity"
	"bot/internal/domain/usecase"
	usecasesteam "bot/internal/domain/usecase/steam"
	reposelenium "bot/internal/repository/selenium"

	"bot/internal/transport/http"

	"github.com/labstack/echo/v4"
)

func Init() {
	e := echo.New()

	steamUser := bot.Init()
	_, err := sqlite.OpenDb()
	if err != nil {
		panic(err)
	}

	seleniumRepo := reposelenium.NewSelenium(steamUser)
	steamUsecase := usecasesteam.NewSteam(seleniumRepo)
	usecase := usecase.NewUseCase(steamUsecase, entity.Bot{SteamUser: steamUser})
	authHandler := http.NewHandler(usecase)
	authHandler.InitRoutes(e)

	e.Logger.Fatal(e.Start(":8088"))
}
