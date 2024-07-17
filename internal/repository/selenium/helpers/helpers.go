package repohelpers

import "github.com/tebeka/selenium"

type IHelpers interface {
	GetLinksForTradeItem(wd selenium.WebDriver, game string) (map[string]float64, error)
}

type helpers struct{}

func NewHelpers() IHelpers {
	return &helpers{}
}
