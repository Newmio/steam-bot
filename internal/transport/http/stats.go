package http

import (
	"fmt"

	"github.com/Newmio/steam_helper"
	"github.com/labstack/echo/v4"
)

func (h *handler) synchItems(c echo.Context) error {
	game := c.QueryParam("game")

	if game == "" || len([]rune(game)) <= 3 {
		return c.JSON(400, steam_helper.Trace(fmt.Errorf("game is empty")).Error())
	}

	if err := h.s.SynchItems(game); err != nil {
		return c.JSON(500, steam_helper.Trace(err).Error())
	}

	return c.JSON(200, "ok")
}
