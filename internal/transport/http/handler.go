package http

import (
	"bot/internal/domain/usecase"
	"fmt"

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
	e.POST("/ping", h.ping)
	e.GET("/login", h.login)
	e.GET("/test", h.test)

	synch := e.Group("/synch")
	{
		synch.GET("/items", h.synchItems)
	}

	check := e.Group("/check")
	{
		items := check.Group("/items")
		{
			items.GET("/steam", h.checkItems)
		}
	}

	helpers := e.Group("/helper")
	{
		trade := helpers.Group("/trade")
		{
			trade.GET("/steam", h.getLinksForTradeItem)
		}
	}
}

func (h *handler) test(c echo.Context) error {
	if err := h.s.GetRareFloats(60, 0, "csgo"); err != nil {
		return c.JSON(500, steam_helper.Trace(err).Error())
	}

	return c.JSON(200, "ok")
}

func (h *handler) ping(c echo.Context) error {
	type Url struct {
		Url string `json:"url"`
	}
	var u Url
	if err := c.Bind(&u); err != nil {
		return c.JSON(400, steam_helper.Trace(err).Error())
	}

	fmt.Println(u.Url)

	html, err := h.s.Ping(u.Url)
	if err != nil {
		return c.JSON(500, steam_helper.Trace(err).Error())
	}

	return c.HTML(200, html)
}
