package entity

type DbSteamSkins struct {
	Id     string `db:"id"`
	Name   string `db:"name"`
	RuName string `db:"runame"`
	Link   string `db:"link"`
}

