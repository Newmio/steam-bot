package entity

type Bot struct {
	SteamUser SteamUser
	//Mu          sync.Mutex
	IsBusy      bool
	Windows     int
	ProfileLink string
}

type Proxy struct {
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Country  string `json:"country"`
	City     string `json:"city"`
}

type SteamUser struct {
	Login    string  `json:"login"`
	Password string  `json:"password"`
	Proxy    []Proxy `json:"proxy"`
}
