package usecasesteam

import (
	"bot/internal/domain/entity"
)

func (u *steam) SteamAuth(user entity.SteamUser) error {
	return u.r.SteamLogin(user)
}
