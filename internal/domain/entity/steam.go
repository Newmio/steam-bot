package entity

type PriceHistoryResponse struct {
	Success bool        `json:"success"`
	Prices  [][]interface{} `json:"prices"`
}

type SteamItemPrice struct {
	DateTime string
	Price    float64
	Quantity int
}