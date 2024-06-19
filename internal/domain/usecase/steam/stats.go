package usecasesteam

import (
	"bot/internal/domain/entity"
	"fmt"
	"time"

	"github.com/Newmio/steam_helper"
)

func (u *steam) GetCSGOSkins(login string) error {
	ch := make(chan interface{})

	go u.r.GetCSGOSkins(login, ch)

	for {
		select {
		case skins := <-ch:

			if skins == nil {
				return nil
			}

			switch v := skins.(type) {
			case []entity.SteamSkin:
				fmt.Println("======================")
				fmt.Printf("%+v\n", v)
				fmt.Println("======================")

			case error:
				return steam_helper.Trace(v)
			}

		case <- time.After(5 * time.Minute):
			return steam_helper.Trace(fmt.Errorf("timeout"))
		}
	}
}
