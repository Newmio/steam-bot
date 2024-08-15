package repohelpers

import (
	"github.com/Newmio/steam_helper"
	"github.com/tebeka/selenium"
)

type IHelpers interface {
	GetLinksForTradeItem(wd selenium.WebDriver, game string) (map[string]float64, error)
}

type helpers struct {
	http steam_helper.ICustomHTTP
}

func NewHelpers(http steam_helper.ICustomHTTP) IHelpers {
	return &helpers{http: http}
}
