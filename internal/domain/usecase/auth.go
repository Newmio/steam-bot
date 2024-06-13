package usecase

import (
	reposelenium "bot/internal/repository/selenium"
)

type IAuth interface {
	SteamAuth(login string) error
}

type auth struct {
	r reposelenium.ISelenium
}

func NewAuth(r reposelenium.ISelenium) IAuth {
	return &auth{r: r}
}

func (u *auth) SteamAuth(login string) error {
	return u.r.SteamLogin(login)
}
