package usecase

import (
	"bot/internal/domain/entity"
	usecasesteam "bot/internal/domain/usecase/steam"
	"fmt"
)

type IUseCase interface {
	SteamAuth(login string) error
}

type useCase struct {
	bots  map[string]entity.Bot
	steam usecasesteam.ISteam
}

func NewUseCase(steam usecasesteam.ISteam, bots map[string]entity.Bot) IUseCase {
	return &useCase{bots: bots, steam: steam}
}

func (u *useCase) SteamAuth(login string) error {
	if bot, ok := u.bots[login]; ok {
		return u.steam.SteamAuth(bot.SteamUser)
	}

	return fmt.Errorf("user %s not found", login)
}
