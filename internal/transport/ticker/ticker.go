package ticker

import "bot/internal/domain/usecase"

type ticker struct {
	s usecase.IUseCase
}

func NewTicker(s usecase.IUseCase) *ticker{
	return &ticker{s: s}
}