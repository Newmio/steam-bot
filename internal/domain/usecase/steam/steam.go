package usecasesteam

import (
	"bot/internal/domain/entity"
	repodb "bot/internal/repository/db"
	reposelenium "bot/internal/repository/selenium"
)

type ISteam interface {
	SteamAuth(user entity.SteamUser) error
	SynchSteamCSGOSkins(login string) error
}

type steam struct {
	r  reposelenium.ISelenium
	db repodb.IDatabase
}

func NewSteam(r reposelenium.ISelenium, db repodb.IDatabase) ISteam {
	return &steam{r: r, db: db}
}
