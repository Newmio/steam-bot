package app

import (
	"bot/internal/configs/bot"
	"bot/internal/configs/sqlite"
	"bot/internal/domain/usecase"
	usecasedmarket "bot/internal/domain/usecase/dmarket"
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
	fmt.Println("------------ 1 ------------")
	e := echo.New()
	fmt.Println("------------ 2 ------------")
	bot := bot.Init()
	fmt.Println("------------ 3 ------------")

	sqlite, err := sqlite.OpenDb(bot.SteamUser.Login)
	if err != nil {
		panic(err)
	}
	fmt.Println("------------ 4 ------------")

	// redis, err := redis.OpenDb()
	// if err != nil {
	// 	panic(err)
	// }

	seleniumRepo := reposelenium.NewSelenium(bot.SteamUser)
	fmt.Println("------------ 5 ------------", seleniumRepo)
	repoRedis := reporedis.NewRedis(nil)
	fmt.Println("------------ 6 ------------")
	repoSqlite := reposqlite.NewSqlite(sqlite)
	fmt.Println("------------ 7 ------------")
	dbRepo := repodb.NewDatabase(repoRedis, repoSqlite)
	fmt.Println("------------ 8 ------------")
	steamUsecase := usecasesteam.NewSteam(seleniumRepo, dbRepo)
	fmt.Println("------------ 9 ------------")
	dmarketUsecase := usecasedmarket.NewDmarket(seleniumRepo, dbRepo)
	fmt.Println("------------ 10 ------------")
	usecase := usecase.NewUseCase(steamUsecase, dmarketUsecase, bot)
	fmt.Println("------------ 11 ------------")
	authHandler := http.NewHandler(usecase)
	fmt.Println("------------ 12 ------------")
	authHandler.InitRoutes(e)
	fmt.Println("------------ 13 ------------")

	e.Logger.Fatal(e.Start(":8088"))
}
