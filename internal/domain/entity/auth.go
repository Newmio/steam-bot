package entity

import "net/http"

type AuthInfo struct {
	Proxy     string
	UserAgent string
	Cookies   []*http.Cookie
}

type ProxyAccount struct{
	SteamLogin string `json:"steam_login"`
	SteamPassword string `json:"steam_password"`
	Ip string `json:"ip"`
	Port string `json:"port"`
	Login string `json:"login"`
	Password string `json:"password"`
	Country string `json:"country"`
	City string `json:"city"`
}