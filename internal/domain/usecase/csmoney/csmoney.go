package usecasecsmoney

import (
	repodb "bot/internal/repository/db"
	reposelenium "bot/internal/repository/selenium"

	"github.com/Newmio/steam_helper"
)

type ICsmoney interface{}

type csmoney struct {
	r  reposelenium.ISelenium
	db repodb.IDatabase
	http steam_helper.ICustomHTTP
}

func NewCsmoney(r reposelenium.ISelenium, db repodb.IDatabase, http steam_helper.ICustomHTTP) ICsmoney {
	return &csmoney{r: r, db: db, http: http}
}