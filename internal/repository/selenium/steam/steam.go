package reposteam

import (
	"bot/internal/domain/entity"

	"github.com/tebeka/selenium"
)

type ISteam interface {
	Login(wd selenium.WebDriver, user entity.User) (string, error)
	GetCSGOStats(wd selenium.WebDriver, ch chan entity.SteamSkin) error
}

type steam struct{}

func NewSteam() ISteam {
	return &steam{}
}
