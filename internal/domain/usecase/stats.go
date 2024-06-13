package usecase

import reposelenium "bot/internal/repository/selenium"

type IStats interface {

}

type stats struct {
	r reposelenium.ISelenium
}

func NewStats(r reposelenium.ISelenium) IStats {
	return &stats{r: r}
}
