package http

import (
	"github.com/Newmio/steam_helper"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Login(c echo.Context) error {
	if err := h.s.SteamAuth(); err != nil {
		return c.HTML(500, steam_helper.Trace(err).Error())
	}

	return c.JSON(200, "ok")
}
