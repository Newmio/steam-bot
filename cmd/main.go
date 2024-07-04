package main

import (
	"bot/internal/app"

	_ "github.com/mattn/go-sqlite3"
)

// chromedriver --port=9515
// docker system prune -a --volumes

func main() {
	app.Init()
}
