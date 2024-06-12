package app

import (
	"bot/internal/configs/account"
	"bot/internal/domain/usecase"
	repoauth "bot/internal/repository/auth"
	"bot/internal/transport/http"

	"github.com/labstack/echo/v4"
)

func Init() {
	e := echo.New()

	authRepo := repoauth.NewAuth(account.InitAccounts())
	authUsecase := usecase.NewAuth(authRepo)
	authHandler := http.NewHandler(authUsecase)
	authHandler.InitRoutes(e)

	e.Logger.Fatal(e.Start(":8088"))
}
