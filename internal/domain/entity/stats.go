package entity

import "time"

type CheckItem struct {
	Sell    []Offer
	Buy     []Offer
	History []SaleHistory
}

type SaleHistory struct {
	Date time.Time
	Cost int
	Count int
}

type Offer struct {
	Cost  int
	Count int
}
