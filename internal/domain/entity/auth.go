package entity

type ProxyAccount struct {
	SteamLogin    string `json:"steam_login"`
	SteamPassword string `json:"steam_password"`
	Ip            string `json:"ip"`
	Port          string `json:"port"`
	Login         string `json:"login"`
	Password      string `json:"password"`
	Country       string `json:"country"`
	City          string `json:"city"`
}

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
