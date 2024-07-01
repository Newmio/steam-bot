package http

import (
	"github.com/Newmio/steam_helper"
	"github.com/labstack/echo/v4"
)

func (h *handler) synchSteamCSGOSkins(c echo.Context) error {
	if err := h.s.SynchSteamCSGOSkins(); err != nil {
		return c.JSON(500, steam_helper.Trace(err).Error())
	}

	return c.JSON(200, "ok")
}
