package repocsmoney

import (
	"bot/internal/domain/entity"
	"net/http"

	"github.com/Newmio/steam_helper"
)

type ICsmoney interface {
	GetRareFloats(limit, offset int) (map[string][]entity.FloatItem, error)
}

type csmoney struct {
	http steam_helper.ICustomHTTP
}

func NewCsmoney(proxy []entity.Proxy) ICsmoney {
	client := &http.Client{}
	var httpProxy []steam_helper.ProxyConfig

	for _, value := range proxy {
		httpProxy = append(httpProxy, steam_helper.ProxyConfig{
			Login: value.Login,
			Pass:  value.Password,
			IP:    value.Ip,
			Port:  value.Port,
		})
	}

	return &csmoney{http: steam_helper.NewHttp(client, httpProxy)}
}
