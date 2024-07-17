package usecasehelpers

import (
	reposelenium "bot/internal/repository/selenium"
	"fmt"

	"github.com/Newmio/steam_helper"
)

type IHelpers interface {
	GetLinksForTradeItem(game string) error
}

type helpers struct {
	r reposelenium.ISelenium
}

func NewHelpers(r reposelenium.ISelenium) IHelpers {
	return &helpers{r: r}
}

func (s *helpers) GetLinksForTradeItem(game string) error {
	links, err := s.r.GetLinksForTradeItem(game)
	if err != nil {
		return steam_helper.Trace(err)
	}

	for key, value := range links {
		fmt.Println(key, ": ", value)
	}

	return nil
}
