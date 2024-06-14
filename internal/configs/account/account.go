package account

import (
	"bot/internal/domain/entity"
	"encoding/json"
	"os"

	"github.com/Newmio/steam_helper"
)

type account struct {
	Accs []entity.SteamUser `json:"accounts"`
}

func InitAccounts() []entity.SteamUser {
	var accounts account

	file, err := os.Open("internal/configs/account/accounts.json")
	if err != nil {
		panic(steam_helper.Trace(err))
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&accounts); err != nil {
		panic(steam_helper.Trace(err))
	}

	return accounts.Accs
}
