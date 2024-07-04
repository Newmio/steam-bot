package bot

import (
	"bot/internal/domain/entity"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Newmio/steam_helper"
)

type botConfig struct{
	Bot entity.Bot `json:"bot"`
	Markets map[string]entity.Market `json:"markets"`
}

func Init() botConfig {
	var botConfig botConfig

	filename := os.Getenv("STEAM_LOGIN")
	if filename == "" {
		filename = "bot"
	}

	file, err := os.Open(fmt.Sprintf("internal/configs/bot/%s.json", filename))
	if err != nil {
		panic(steam_helper.Trace(err))
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&botConfig); err != nil {
		panic(steam_helper.Trace(err))
	}

	return botConfig
}
