package reposqlite

import (
	"bot/internal/domain/entity"

	"github.com/Newmio/steam_helper"
	"github.com/jmoiron/sqlx"
)

type ISqlite interface {
	CreateTables() error
	CreateSteamSkins(skins []entity.DbSteamSkins) error
	CreateStickerSkins(skins []entity.DbSteamSkins) error
	CreateFloatSkins(skins []entity.DbSteamSkins) error
	CreatePatternSkins(skins []entity.DbSteamSkins) error
	CreateBetweenSkins(skins []entity.DbSteamSkins) error
	GetSteamSkins(limit, offset int) ([]entity.DbSteamSkins, error)
	GetStickerSkins(limit, offset int) ([]entity.DbSteamSkins, error)
	GetFloatSkins(limit, offset int) ([]entity.DbSteamSkins, error)
	GetPatternSkins(limit, offset int) ([]entity.DbSteamSkins, error)
	GetBetweenSkins(limit, offset int) ([]entity.DbSteamSkins, error)
}

type sqlite struct {
	db *sqlx.DB
}

func NewSqlite(db *sqlx.DB) ISqlite {
	r := &sqlite{db: db}

	if err := r.CreateTables(); err != nil {
		panic(steam_helper.Trace(err))
	}

	return r
}

func (db *sqlite) CreateTables() error {
	str := `create table if not exists steam_skins(
		id text primary key,
		name text default '',
		runame text default '',
		link text default ''
	)`

	if _, err := db.db.Exec(str); err != nil {
		return steam_helper.Trace(err)
	}

	str = `create table if not exists sticker_skins(
		id text primary key,
		name text default '',
		runame text default '',
		link text default ''
	)`

	if _, err := db.db.Exec(str); err != nil {
		return steam_helper.Trace(err)
	}

	str = `create table if not exists float_skins(
		id text primary key,
		name text default '',
		runame text default '',
		link text default ''
	)`

	if _, err := db.db.Exec(str); err != nil {
		return steam_helper.Trace(err)
	}

	str = `create table if not exists pattern_skins(
		id text primary key,
		name text default '',
		runame text default '',
		link text default ''
	)`

	if _, err := db.db.Exec(str); err != nil {
		return steam_helper.Trace(err)
	}

	str = `create table if not exists between_skins(
		id text primary key,
		csmoney boolean default false,
		dmarket boolean default false,
		name text default '',
		runame text default '',
		link text default ''
	)`

	if _, err := db.db.Exec(str); err != nil {
		return steam_helper.Trace(err)
	}

	return nil
}
