package http

import (
	"github.com/Newmio/steam_helper"
	"github.com/labstack/echo/v4"
)

func (h *Handler) SynchSteamCSGOSkins(c echo.Context) error {

	if err := h.s.SynchSteamCSGOSkins(c.QueryParam("login")); err != nil {
		return c.HTML(500, steam_helper.Trace(err).Error())
	}

	return c.JSON(200, "ok")
}
