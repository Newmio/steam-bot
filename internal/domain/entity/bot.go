package entity

import (
	"time"
)

type Bot struct {
	SteamUser SteamUser `json:"steam_user"`
	//Mu          sync.Mutex
	Markets       map[string]Market `json:"markets"`
	DateStopBot   time.Time         `json:"date_stop_bot"`
	MaxSeleniumWd int               `json:"max_selenium_wd"`
	IsBusy        bool
	Windows       int
}

type Market struct {
	Trade        trade        `json:"trade"`
	StickerTrade stickerTrade `json:"sticker_trade"`
	FloatTrade   floatTrade   `json:"float_trade"`
	PatternTrade patternTrade `json:"pattern_trade"`
}

type trade struct {
	Do           bool         `json:"do"`
	RangeDay     int          `json:"range_day"`
	MinCount     int          `json:"min_count"`
	MinSellCount int          `json:"min_sell_count"`
	MinCost      float64      `json:"min_cost"`
	MinSellCost  float64      `json:"min_sell_cost"`
	Offers       []tradeOffer `json:"offers"`
}

type stickerTrade struct {
	Do             bool            `json:"do"`
	MinProfit      int             `json:"min_profit"`
	LikedStickers  []likedSticker  `json:"liked_stickers"`
	IgnoreStickers []ignoreSticker `json:"ignore_stickers"`
}

type floatTrade struct{}

type patternTrade struct{}

type likedSticker struct {
	ItemName    string  `json:"item_name"`
	ItemHash    string  `json:"item_hash"`
	StickerName string  `json:"sticker_name"`
	StickerHash string  `json:"sticker_hash"`
	Cost        float64 `json:"cost"`
	One         bool    `json:"1"`
	Two         bool    `json:"2"`
	Three       bool    `json:"3"`
	Four        bool    `json:"4"`
}

type ignoreSticker struct {
	StickerName string `json:"sticker_name"`
	StickerHash string `json:"sticker_hash"`
}

type tradeOffer struct {
	MinProfit int `json:"min_profit"`
	MaxProfit int `json:"max_profit"`
	Count     int `json:"count"`
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
		//market := bot.Markets[marketName]

		switch action {
		default:
			return true
		}
	}

	return false
}
