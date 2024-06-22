package http

import (
	"bot/internal/domain/usecase"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	s usecase.IUseCase
}

func NewHandler(s usecase.IUseCase) *Handler {
	return &Handler{s: s}
}

func (h *Handler) InitRoutes(e *echo.Echo) {
	e.GET("/login", h.Login)

	synch := e.Group("/synch")
	{
		csgo := synch.Group("/csgo")
		{
			csgo.GET("/steam", h.SynchSteamCSGOSkins)
		}
	}
}
