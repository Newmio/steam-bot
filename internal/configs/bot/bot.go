package bot

import (
	"bot/internal/domain/entity"
	"encoding/json"
	"os"

	"github.com/Newmio/steam_helper"
)

func Init() entity.Bot {
	var bot entity.Bot

	file, err := os.Open("internal/configs/bot/bot.json")
	if err != nil {
		panic(steam_helper.Trace(err))
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&bot); err != nil {
		panic(steam_helper.Trace(err))
	}

	return bot
}
