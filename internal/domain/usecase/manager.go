package usecase

import (
	"bot/internal/domain/entity"
	usecasesteam "bot/internal/domain/usecase/steam"
)

type IUseCase interface {
	SteamAuth(login string) error
	SynchSteamCSGOSkins(login string) error
}

type useCase struct {
	bot   entity.Bot
	steam usecasesteam.ISteam
}

func NewUseCase(steam usecasesteam.ISteam, bot entity.Bot) IUseCase {
	return &useCase{bot: bot, steam: steam}
}

func (u *useCase) SynchSteamCSGOSkins(login string) error {
	if !u.bot.CheckAction("synch"){
		return nil
	}
	return u.steam.SynchSteamCSGOSkins(u.bot.SteamUser.Login)
}

func (u *useCase) SteamAuth(login string) error {
	return u.steam.SteamAuth(u.bot.SteamUser)
}
