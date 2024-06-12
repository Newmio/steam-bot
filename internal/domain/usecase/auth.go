package usecase

import (
	"bot/internal/domain/entity"
	repoauth "bot/internal/repository/auth"
)

type IAuth interface {
	Login(login string) (entity.AuthInfo, error)
}

type auth struct {
	r repoauth.IAuth
}

func NewAuth(r repoauth.IAuth) IAuth {
	return &auth{r: r}
}

func (u *auth) Login(login string) (entity.AuthInfo, error) {
	return u.r.Login(login)
}
