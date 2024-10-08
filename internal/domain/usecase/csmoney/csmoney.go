package usecasecsmoney

import (
	repodb "bot/internal/repository/db"
	reposelenium "bot/internal/repository/selenium"
)

type ICsmoney interface{
	GetRareFloats(limit, offset int, game string) error
}

type csmoney struct {
	r  reposelenium.ISelenium
	db repodb.IDatabase
}

func NewCsmoney(r reposelenium.ISelenium, db repodb.IDatabase) ICsmoney {
	return &csmoney{r: r, db: db}
}
