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
	e.GET("/login", h.Login)
}

func (h *Handler) Login(c echo.Context) error {
	if err := h.s.SteamAuth(c.QueryParam("login")); err != nil {
		return c.HTML(500, steam_helper.Trace(err).Error())
	}

	return c.JSON(200, "ok")
}
