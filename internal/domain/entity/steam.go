package entity

import "time"

type SteamSellHistory struct {
	HashName    string
	PriceSuffix string
	Price      SteamItemPrice
}

type SteamItemPrice struct {
	DateTime time.Time
	Cost     int
	Count    int
}
