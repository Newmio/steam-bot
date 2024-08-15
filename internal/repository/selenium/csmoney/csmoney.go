package repocsmoney

import (
	"bot/internal/domain/entity"

	"github.com/Newmio/steam_helper"
)

type ICsmoney interface {
	GetRareFloats(limit, offset int) (map[string][]entity.FloatItem, error)
}

type csmoney struct {
	http steam_helper.ICustomHTTP
}

func NewCsmoney(http steam_helper.ICustomHTTP) ICsmoney {
	return &csmoney{http: http}
}
