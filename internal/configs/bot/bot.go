package bot

import (
	"bot/internal/domain/entity"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Newmio/steam_helper"
)

func Init() entity.Bot {
	var bot entity.Bot

	filename := os.Getenv("STEAM_LOGIN")
	if filename == "" {
		filename = "bot"
	}

	file, err := os.Open(fmt.Sprintf("internal/configs/bot/%s.json", filename))
	if err != nil {
		panic(steam_helper.Trace(err))
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&bot); err != nil {
		panic(steam_helper.Trace(err))
	}

	return bot
}
