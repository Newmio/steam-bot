package proxy

import (
	"bot/internal/domain/entity"
	"encoding/json"
	"os"

	"github.com/Newmio/steam_helper"
)

type proxy struct {
	Proxy []entity.Proxy `json:"proxy"`
}

func InitProxy() []entity.Proxy {
	var proxy proxy

	file, err := os.Open("internal/configs/account/accounts.json")
	if err != nil {
		panic(steam_helper.Trace(err))
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&proxy); err != nil {
		panic(steam_helper.Trace(err))
	}

	return proxy.Proxy
}
