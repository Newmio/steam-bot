package entity

import (
	"time"
)

type Bot struct {
	SteamUser SteamUser `json:"steam_user"`
	//Mu          sync.Mutex
	Markets       map[string]market `json:"markets"`
	DateStopBot   time.Time         `json:"date_stop_bot"`
	MaxSeleniumWd int               `json:"max_selenium_wd"`
	IsBusy        bool
	Windows       int
}

type market struct {
	MinCount     int     `json:"min_count"`
	Trade        bool    `json:"trade"`
	StickerTrade bool    `json:"sticker_trade"`
	FloatTrade   bool    `json:"float_trade"`
	PatternTrade bool    `json:"pattern_trade"`
}

type proxy struct {
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Country  string `json:"country"`
	City     string `json:"city"`
}

type SteamUser struct {
	Id          string  `json:"id"`
	ProfileLink string  `json:"profile_link"`
	Login       string  `json:"login"`
	Password    string  `json:"password"`
	Proxy       []proxy `json:"proxy"`
}

func (bot *Bot) CheckAction(marketName, action string) bool {

	if bot.DateStopBot.After(time.Now()) && !bot.IsBusy {
		market := bot.Markets[marketName]

		switch action {
		case "trade":
			return market.Trade
		case "sticker_trade":
			return market.StickerTrade
		case "float_trade":
			return market.FloatTrade
		case "pattern_trade":
			return market.PatternTrade
		default:
			return true
		}
	}

	return false
}
