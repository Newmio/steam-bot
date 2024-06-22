package entity

import "time"

type Bot struct {
	SteamUser SteamUser `json:"steam_user"`
	//Mu          sync.Mutex
	MinSynchCost  float32   `json:"min_synch_cost"`
	MaxSynchCost  float32   `json:"max_synch_cost"`
	Synch         bool      `json:"synch"`
	SteamTrade    bool      `json:"steam_trade"`
	DmarketTrade  bool      `json:"dmarket_trade"`
	CsmoneyTrade  bool      `json:"csmoney_trade"`
	StickerTrade  bool      `json:"sticker_trade"`
	FloatTrade    bool      `json:"float_trade"`
	PatternTrade  bool      `json:"pattern_trade"`
	DateStopBot   time.Time `json:"date_stop_bot"`
	MaxSeleniumWd int       `json:"max_selenium_wd"`
	IsBusy        bool
	Windows       int
}

type Proxy struct {
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
	Proxy       []Proxy `json:"proxy"`
}

func (bot *Bot) CheckAction(action string) bool {

	if bot.DateStopBot.After(time.Now()) && !bot.IsBusy {
		switch action {
		case "synch":
			return bot.Synch
		case "steam_trade":
			return bot.SteamTrade
		case "dmarket_trade":
			return bot.DmarketTrade
		case "csmoney_trade":
			return bot.CsmoneyTrade
		case "sticker_trade":
			return bot.StickerTrade
		case "float_trade":
			return bot.FloatTrade
		case "pattern_trade":
			return bot.PatternTrade
		}
	}

	return false
}
