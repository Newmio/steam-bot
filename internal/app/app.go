package app

import (
	"bot/internal/configs/account"
	"bot/internal/domain/usecase"
	reposelenium "bot/internal/repository/selenium"

	"bot/internal/transport/http"

	"github.com/labstack/echo/v4"
)

func Init() {
	e := echo.New()

	authRepo := reposelenium.NewSelenium(account.InitAccounts())
	authUsecase := usecase.NewAuth(authRepo)
	authHandler := http.NewHandler(authUsecase)
	authHandler.InitRoutes(e)

	e.Logger.Fatal(e.Start(":8088"))
}
