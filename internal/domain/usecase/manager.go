package usecase

import (
	"bot/internal/domain/entity"
	usecasedmarket "bot/internal/domain/usecase/dmarket"
	usecasesteam "bot/internal/domain/usecase/steam"
	"sync"

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
	bot.Wg = new(sync.WaitGroup)
	return &useCase{bot: bot, steam: steam, dmarket: dmarket}
}

func (s *useCase) CheckTradeItems(game string, start, stop int) error {
	s.bot.Wg.Wait()
	s.bot.Wg.Add(1)
	defer s.bot.Wg.Done()

	return s.steam.CheckTradeItems(game, start, stop)
}

func (s *useCase) SynchItems(game string) error {
	if !s.bot.Synch {
		return nil
	}
	
	s.bot.Wg.Wait()
	s.bot.Wg.Add(1)
	defer s.bot.Wg.Done()

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
	s.bot.Wg.Wait()
	s.bot.Wg.Add(1)
	defer s.bot.Wg.Done()

	return s.steam.SteamAuth()
}

func (s *useCase) Ping(url string) (string, error) {
	return s.steam.Ping(url)
}
