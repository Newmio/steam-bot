package usecasesteam

import (
	"bot/internal/domain/entity"
	reposelenium "bot/internal/repository/selenium"
)

type ISteam interface {
	SteamAuth(user entity.SteamUser) error
	GetCSGOSkins(login string) error
}

type steam struct {
	r reposelenium.ISelenium
}

func NewSteam(r reposelenium.ISelenium) ISteam {
	return &steam{r: r}
}
