package http

import (
	"bot/internal/domain/usecase"

	"github.com/Newmio/steam_helper"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	s usecase.IAuth
}

func NewHandler(s usecase.IAuth) *Handler {
	return &Handler{s: s}
}

func (h *Handler) InitRoutes(e *echo.Echo) {
	e.GET("/login", h.Login)
}

func (h *Handler) Login(c echo.Context) error {
	login := "jasgt26800"

	_, err := h.s.Login(login)
	if err != nil {
		return c.JSON(500, steam_helper.Trace(err).Error())
	}

	return c.JSON(200, "ok")
}
