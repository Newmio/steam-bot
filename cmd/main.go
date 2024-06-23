package main

import (
	"bot/internal/app"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// chromedriver --port=9515

func main() {
	os.Setenv("MOZ_REMOTE_SETTINGS_DEVTOOLS", "1")
	app.Init()
}
