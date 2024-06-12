package account

import (
	"bot/internal/domain/entity"
	"encoding/json"
	"os"

	"github.com/Newmio/steam_helper"
)

type proxy struct {
	P []entity.ProxyAccount `json:"proxy"`
}

func InitAccounts() map[string][]entity.ProxyAccount {
	accounts := make(map[string][]entity.ProxyAccount)
	var proxies proxy

	file, err := os.Open("internal/configs/account/proxy.json")
	if err != nil {
		panic(steam_helper.Trace(err))
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&proxies); err != nil {
		panic(steam_helper.Trace(err))
	}

	for _, value := range proxies.P {
		if value.SteamLogin != "" && value.SteamPassword != "" {
			accounts[value.SteamLogin] = append(accounts[value.SteamLogin], value)
		}
	}

	return accounts
}
