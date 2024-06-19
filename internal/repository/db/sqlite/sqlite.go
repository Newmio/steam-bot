package reposqlite

import (
	"github.com/Newmio/steam_helper"
	"github.com/jmoiron/sqlx"
)

type ISqlite interface{}

type sqlite struct {
	db *sqlx.DB
}

func NewSqlite(db *sqlx.DB) ISqlite {
	return &sqlite{db: db}
}

func (db *sqlite) CreateTables() error {
	str := `create table if not exists steam_skins(
		id text primary key,
		name text,
		runame text)`

	if _, err := db.db.Exec(str); err != nil {
		return steam_helper.Trace(err)
	}

	str = `create table if not exists sticker_skins(
		id text primary key,
		name text,
		runame text
	)`

	if _, err := db.db.Exec(str); err != nil {
		return steam_helper.Trace(err)
	}

	str = `create table if not exists float_skins(
		id text primary key,
		name text,
		runame text
	)`

	if _, err := db.db.Exec(str); err != nil {
		return steam_helper.Trace(err)
	}

	str = `create table if not exists pattern_skins(
		id text primary key,
		name text,
		runame text
	)`

	if _, err := db.db.Exec(str); err != nil {
		return steam_helper.Trace(err)
	}

	return nil
}
