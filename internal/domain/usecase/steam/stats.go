package usecasesteam

import (
	"bot/internal/domain/entity"
	"fmt"
)

func (u *steam) GetCSGOSkins(login string) error {
	ch := make(chan []entity.SteamSkin)

	go u.r.GetCSGOSkins(login, ch)

	fmt.Println("======================")
	fmt.Printf("%+v\n", <-ch)
	fmt.Println("======================")

	return nil
}
