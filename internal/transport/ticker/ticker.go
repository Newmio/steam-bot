package ticker

import "bot/internal/domain/usecase"

type ticker struct {
	s usecase.IUseCase
	games []string
}

func NewTicker(s usecase.IUseCase)*ticker{
	return &ticker{
		s: s,
		games: []string{"csgo", "dota2"},
	}
}

func (t *ticker) synchItems()error{
	return nil
}