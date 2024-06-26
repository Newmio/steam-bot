package usecasedmarket

import (
	repodb "bot/internal/repository/db"
	reposelenium "bot/internal/repository/selenium"
)

type IDmarket interface{}

type dmarket struct {
	r  reposelenium.ISelenium
	db repodb.IDatabase
}

func NewDmarket(r reposelenium.ISelenium, db repodb.IDatabase) IDmarket {
	return &dmarket{r: r, db: db}
}
