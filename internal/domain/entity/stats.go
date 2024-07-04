package entity

import "time"

type CheckItem struct {
	Sell    []CheckOffer
	Buy     []CheckOffer
	History []SaleHistory
}

type SaleHistory struct {
	Date time.Time
	Cost int
	Count int
}

type CheckOffer struct {
	Cost  int
	Count int
}
