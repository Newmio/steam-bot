package usecasecsmoney

import (
	"fmt"

	"github.com/Newmio/steam_helper"
)

func (s *csmoney) GetRareFloats(limit, offset int, game string) error {
	items, err := s.r.GetRareFloats(limit, offset)
	if err != nil {
		return steam_helper.Trace(err)
	}

	if err := s.db.CreateItemsRareFloat(items, game); err != nil {
		return steam_helper.Trace(err)
	}

	fmt.Println("---------------------------------------------")
	for key, value := range items {
		fmt.Println("+++++++++++++++++++++++++++++++++++++++++")
		for _, item := range value {
			fmt.Printf("%s: float = %f cost = %d\n", key, item.Float, item.Cost)
		}
		fmt.Println("+++++++++++++++++++++++++++++++++++++++++")
	}
	fmt.Println("---------------------------------------------")

	return nil
}
