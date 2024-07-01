package usecasesteam

import (
	repodb "bot/internal/repository/db"
	reposelenium "bot/internal/repository/selenium"
)

type ISteam interface {
	SteamAuth() error
	SynchCSGOSkins(minCost, maxCost float64, minCount int) error
	Ping(url string) (string, error)
}

type steam struct {
	r  reposelenium.ISelenium
	db repodb.IDatabase
}

func NewSteam(r reposelenium.ISelenium, db repodb.IDatabase) ISteam {
	return &steam{r: r, db: db}
}

func (s *steam) Ping(url string) (string, error) {
	return s.r.Ping(url)
}