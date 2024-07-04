package usecase

import (
	"bot/internal/domain/entity"
	usecasedmarket "bot/internal/domain/usecase/dmarket"
	usecasesteam "bot/internal/domain/usecase/steam"
)

type IUseCase interface {
	SteamAuth() error
	SynchItems(game string) error
	Ping(url string) (string, error)
}

type useCase struct {
	bot     entity.Bot
	steam   usecasesteam.ISteam
	dmarket usecasedmarket.IDmarket
}

func NewUseCase(steam usecasesteam.ISteam, dmarket usecasedmarket.IDmarket, bot entity.Bot) IUseCase {
	return &useCase{bot: bot, steam: steam, dmarket: dmarket}
}

func (u *useCase) SynchItems(game string) error {
	return u.steam.SynchItems(game)
}

func (u *useCase) SteamAuth() error {
	if u.bot.CheckAction("", "") {
		u.bot.IsBusy = true
		return u.steam.SteamAuth()
	}

	return nil
}

func (u *useCase) Ping(url string) (string, error) {
	return u.steam.Ping(url)
}
