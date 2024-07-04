package http

import (
	"bot/internal/domain/usecase"

	"github.com/Newmio/steam_helper"
	"github.com/labstack/echo/v4"
)

type handler struct {
	s usecase.IUseCase
}

func NewHandler(s usecase.IUseCase) *handler {
	return &handler{s: s}
}

func (h *handler) InitRoutes(e *echo.Echo) {
	e.GET("/ping", h.ping)
	e.GET("/login", h.login)

	synch := e.Group("/synch")
	{
		synch.GET("/csgo", h.synchItems)
	}
}

func (h *handler) ping(c echo.Context) error {
	html, err := h.s.Ping("https://" + c.QueryParam("url"))
	if err != nil {
		return c.JSON(500, steam_helper.Trace(err).Error())
	}

	return c.HTML(200, html)
}