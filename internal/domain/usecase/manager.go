package usecase

import (
	"bot/internal/domain/entity"
	usecasedmarket "bot/internal/domain/usecase/dmarket"
	usecasesteam "bot/internal/domain/usecase/steam"
)

type IUseCase interface {
	SteamAuth() error
	SynchSteamCSGOSkins() error
	SynchDmarketCSGOSkins() error
	Ping(url string) error
}

type useCase struct {
	bot     entity.Bot
	steam   usecasesteam.ISteam
	dmarket usecasedmarket.IDmarket
}

func NewUseCase(steam usecasesteam.ISteam, dmarket usecasedmarket.IDmarket, bot entity.Bot) IUseCase {
	return &useCase{bot: bot, steam: steam, dmarket: dmarket}
}

func (u *useCase) SynchDmarketCSGOSkins() error {
	market := u.bot.Markets["dmarket"]

	if !u.bot.CheckAction("dmarket", "synch") {
		return nil
	}
	u.bot.IsBusy = true
	return u.steam.SynchCSGOSkins(
		market.MinSynchCost,
		market.MaxSynchCost,
		market.MinCount,
	)
}

func (u *useCase) SynchSteamCSGOSkins() error {
	market := u.bot.Markets["dmarket"]

	if !u.bot.CheckAction("steam", "synch") {
		return nil
	}
	u.bot.IsBusy = true
	return u.steam.SynchCSGOSkins(
		market.MinSynchCost,
		market.MaxSynchCost,
		market.MinCount,
	)
}

func (u *useCase) SteamAuth() error {
	if !u.bot.CheckAction("", "") {
		return nil
	}
	u.bot.IsBusy = true
	return u.steam.SteamAuth()
}

func (u *useCase) Ping(url string) error{
	return u.steam.Ping(url)
}