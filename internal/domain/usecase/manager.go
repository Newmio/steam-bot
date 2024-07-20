package usecase

import (
	"bot/internal/domain/entity"
	usecasedmarket "bot/internal/domain/usecase/dmarket"
	usecasehelpers "bot/internal/domain/usecase/helpers"
	usecasesteam "bot/internal/domain/usecase/steam"

	"github.com/Newmio/steam_helper"
)

type IUseCase interface {
	SteamAuth() error
	SynchItems(game string) error
	Ping(url string) (string, error)
	CheckTradeItems(game string, start, stop int) error
	GetLinksForTradeItem(game string) error
}

type useCase struct {
	bot     entity.Bot
	steam   usecasesteam.ISteam
	dmarket usecasedmarket.IDmarket
	helpers usecasehelpers.IHelpers
}

func NewUseCase(steam usecasesteam.ISteam, dmarket usecasedmarket.IDmarket, helpers usecasehelpers.IHelpers, bot entity.Bot) IUseCase {
	return &useCase{bot: bot, steam: steam, dmarket: dmarket, helpers: helpers}
}

func (s *useCase) GetLinksForTradeItem(game string) error {
	return s.helpers.GetLinksForTradeItem(game)
}

func (s *useCase) CheckTradeItems(game string, start, stop int) error {
	return s.steam.CheckTradeItems(game, start, stop)
}

func (s *useCase) SynchItems(game string) error {
	if !s.bot.Synch {
		return nil
	}

	ch := make(steam_helper.CursorCh[[]entity.SteamItem])
	info := entity.PaginationInfo[[]entity.SteamItem]{
		Game:  game,
		Start: s.bot.SynchStart,
		Stop:  s.bot.SynchStop,
		Ch:    ch,
	}

	return s.steam.SynchItems(info)
}

func (s *useCase) SteamAuth() error {
	return s.steam.SteamAuth()
}

func (s *useCase) Ping(url string) (string, error) {
	return s.steam.Ping(url)
}
