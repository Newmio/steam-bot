package reposqlite

import (
	"bot/internal/domain/entity"

	"github.com/Newmio/steam_helper"
)

func (db *sqlite) GetBetweenSkins(limit, offset int) ([]entity.DbSteamSkins, error) {
	var skins []entity.DbSteamSkins

	str := "select * from between_skins order by id limit ? offset ?"

	if err := db.db.Select(&skins, str, limit, offset); err != nil {
		return nil, steam_helper.Trace(err, str)
	}

	return skins, nil
}

func (db *sqlite) GetPatternSkins(limit, offset int) ([]entity.DbSteamSkins, error) {
	var skins []entity.DbSteamSkins

	str := "select * from pattern_skins order by id limit ? offset ?"

	if err := db.db.Select(&skins, str, limit, offset); err != nil {
		return nil, steam_helper.Trace(err, str)
	}

	return skins, nil
}

func (db *sqlite) CreatePatternSkins(skins []entity.DbSteamSkins) error {
	str := `insert or replace into pattern_skins(id, name, runame, link) 
	values(?, ?, ?, ?)`

	stmt, err := db.db.Preparex(str)
	if err != nil {
		return steam_helper.Trace(err, str)
	}

	for _, value := range skins {
		_, err := stmt.Exec(str, value.Id, value.Name, value.RuName, value.Link)
		if err != nil {
			return steam_helper.Trace(err, str)
		}
	}

	return nil
}

func (db *sqlite) GetFloatSkins(limit, offset int) ([]entity.DbSteamSkins, error) {
	var skins []entity.DbSteamSkins

	str := "select * from float_skins order by id limit ? offset ?"

	if err := db.db.Select(&skins, str, limit, offset); err != nil {
		return nil, steam_helper.Trace(err, str)
	}

	return skins, nil
}

func (db *sqlite) CreateFloatSkins(skins []entity.DbSteamSkins) error {
	str := `insert or replace into float_skins(id, name, runame, link) 
	values(?, ?, ?, ?)`

	stmt, err := db.db.Preparex(str)
	if err != nil {
		return steam_helper.Trace(err, str)
	}

	for _, value := range skins {
		_, err := stmt.Exec(str, value.Id, value.Name, value.RuName, value.Link)
		if err != nil {
			return steam_helper.Trace(err, str)
		}
	}

	return nil
}

func (db *sqlite) GetStickerSkins(limit, offset int) ([]entity.DbSteamSkins, error) {
	var skins []entity.DbSteamSkins

	str := "select * from sticker_skins order by id limit ? offset ?"

	if err := db.db.Select(&skins, str, limit, offset); err != nil {
		return nil, steam_helper.Trace(err, str)
	}

	return skins, nil
}

func (db *sqlite) CreateStickerSkins(skins []entity.DbSteamSkins) error {
	str := `insert or replace into sticker_skins(id, name, runame, link) 
	values(?, ?, ?, ?)`

	stmt, err := db.db.Preparex(str)
	if err != nil {
		return steam_helper.Trace(err, str)
	}

	for _, value := range skins {
		_, err := stmt.Exec(str, value.Id, value.Name, value.RuName, value.Link)
		if err != nil {
			return steam_helper.Trace(err, str)
		}
	}

	return nil
}

func (db *sqlite) GetSteamSkins(limit, offset int) ([]entity.DbSteamSkins, error) {
	var skins []entity.DbSteamSkins

	str := "select * from steam_skins order by id limit ? offset ?"

	if err := db.db.Select(&skins, str, limit, offset); err != nil {
		return nil, steam_helper.Trace(err, str)
	}

	return skins, nil
}

func (db *sqlite) CreateSteamSkins(skins []entity.DbSteamSkins) error {
	str := `insert or replace into steam_skins(id, name, runame, link) 
	values(?, ?, ?, ?)`

	stmt, err := db.db.Preparex(str)
	if err != nil {
		return steam_helper.Trace(err, str)
	}

	for _, value := range skins {
		_, err := stmt.Exec(str, value.Id, value.Name, value.RuName, value.Link)
		if err != nil {
			return steam_helper.Trace(err, str)
		}
	}

	return nil
}

func (db *sqlite) CreateBetweenSkins(skins []entity.DbSteamSkins) error {
	str := `insert or replace into between_skins(id, name, runame, link) 
	values(?, ?, ?, ?)`

	stmt, err := db.db.Preparex(str)
	if err != nil {
		return steam_helper.Trace(err, str)
	}

	for _, value := range skins {
		_, err := stmt.Exec(str, value.Id, value.Name, value.RuName, value.Link)
		if err != nil {
			return steam_helper.Trace(err, str)
		}
	}

	return nil
}
