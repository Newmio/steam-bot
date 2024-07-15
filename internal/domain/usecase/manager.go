package usecase

import (
	"bot/internal/domain/entity"
	usecasedmarket "bot/internal/domain/usecase/dmarket"
	usecasesteam "bot/internal/domain/usecase/steam"

	"github.com/Newmio/steam_helper"
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
	ch := make(steam_helper.CursorCh[[]entity.SteamItem])
	info := entity.PaginationInfo[[]entity.SteamItem]{
		Game:  game,
		Start: 1,
		Stop:  100,
		Ch:    ch,
	}
	return s.steam.SynchItems(info)
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
