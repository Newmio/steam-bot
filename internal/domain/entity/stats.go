package entity

import (
	"time"
)

type CheckItem struct {
	HashName string
	Floats   []FloatItem
	Sell     map[int]int
	Buy      map[int]int
}

type SaleHistory struct {
	Date  time.Time
	Cost  int
	Count int
}
