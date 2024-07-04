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
	CheckTradeItems(game string, start, stop int) error
}

type useCase struct {
	bot     entity.Bot
	steam   usecasesteam.ISteam
	dmarket usecasedmarket.IDmarket
}

func NewUseCase(steam usecasesteam.ISteam, dmarket usecasedmarket.IDmarket, bot entity.Bot) IUseCase {
	return &useCase{bot: bot, steam: steam, dmarket: dmarket}
}

func (s *useCase) CheckTradeItems(game string, start, stop int) error {
	return s.steam.CheckTradeItems(game, start, stop)
}

func (s *useCase) SynchItems(game string) error {
	return s.steam.SynchItems(game)
}

func (s *useCase) SteamAuth() error {
	if s.bot.CheckAction("", "") {
		s.bot.IsBusy = true
		return s.steam.SteamAuth()
	}

	s.bot.IsBusy = false

	return nil
}

func (s *useCase) Ping(url string) (string, error) {
	return s.steam.Ping(url)
}
