package app

import (
	"bot/internal/configs/bot"
	"bot/internal/configs/redis"
	"bot/internal/configs/sqlite"
	"bot/internal/domain/usecase"
	usecasesteam "bot/internal/domain/usecase/steam"
	repodb "bot/internal/repository/db"
	reporedis "bot/internal/repository/db/redis"
	reposqlite "bot/internal/repository/db/sqlite"
	reposelenium "bot/internal/repository/selenium"
	"fmt"

	"bot/internal/transport/http"

	"github.com/labstack/echo/v4"
)

func Init() {
	e := echo.New()
	bot := bot.Init()

	fmt.Println(bot)
	fmt.Println("======")

	sqlite, err := sqlite.OpenDb()
	if err != nil {
		panic(err)
	}

	redis, err := redis.OpenDb()
	if err != nil {
		panic(err)
	}

	seleniumRepo := reposelenium.NewSelenium(bot.SteamUser)
	repoRedis := reporedis.NewRedis(redis)
	repoSqlite := reposqlite.NewSqlite(sqlite)
	dbRepo := repodb.NewDatabase(repoRedis, repoSqlite)
	steamUsecase := usecasesteam.NewSteam(seleniumRepo, dbRepo)
	usecase := usecase.NewUseCase(steamUsecase, bot)
	authHandler := http.NewHandler(usecase)
	authHandler.InitRoutes(e)

	e.Logger.Fatal(e.Start(":8088"))
}
