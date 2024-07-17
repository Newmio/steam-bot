package http

import (
	"fmt"
	"strconv"

	"github.com/Newmio/steam_helper"
	"github.com/labstack/echo/v4"
)

func (h *handler) checkTradeItems(c echo.Context) error {
	game := c.QueryParam("game")

	if game == "" || len([]rune(game)) <= 3 {
		return c.JSON(400, steam_helper.Trace(fmt.Errorf("game is empty")).Error())
	}

	start, err := strconv.Atoi(c.QueryParam("start"))
	if err != nil {
		return c.JSON(400, steam_helper.Trace(err).Error())
	}

	stop, err := strconv.Atoi(c.QueryParam("stop"))
	if err != nil {
		return c.JSON(400, steam_helper.Trace(err).Error())
	}

	if err := h.s.CheckTradeItems(game, start, stop); err != nil {
		return c.JSON(500, steam_helper.Trace(err).Error())
	}

	return c.JSON(200, "ok")
}

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
