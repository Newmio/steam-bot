package http

import (
	"bot/internal/domain/usecase"

	"github.com/Newmio/steam_helper"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	s usecase.IUseCase
}

func NewHandler(s usecase.IUseCase) *Handler {
	return &Handler{s: s}
}

func (h *Handler) InitRoutes(e *echo.Echo) {
	e.GET("/ping", h.Ping)
	e.GET("/login", h.Login)

	synch := e.Group("/synch")
	{
		csgo := synch.Group("/csgo")
		{
			csgo.GET("/steam", h.SynchSteamCSGOSkins)
		}
	}
}

func (h *Handler) Ping(c echo.Context) error {
	if err := h.s.Ping(c.QueryParam("url")); err != nil {
		return c.HTML(500, steam_helper.Trace(err).Error())
	}

	return c.JSON(200, "ok")
}
