package usecasesteam

import (
	"bot/internal/domain/entity"
	"fmt"
	"time"

	"github.com/Newmio/steam_helper"
)

func (u *steam) SynchSteamCSGOSkins(login string) error {
	ch := make(steam_helper.CursorCh[[]entity.SeleniumSteamSkin])

	go u.r.SynchSteamCSGOSkins(login, ch)

	for {
		select{

		case skin := <- ch:
			if skin.Error != nil{
				return steam_helper.Trace(skin.Error)
			}

			if err := u.db.CreateSeleniumSteamSkins(skin.Model); err != nil{
				return steam_helper.Trace(err)
			} 

		case <- time.After(5 * time.Minute):
			return steam_helper.Trace(fmt.Errorf("timeout"))
		}
	}
}
