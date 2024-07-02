package entity

type ItemInfo struct{
	Sell []Offer
	Buy []Offer
	History []SaleHistory
	//Available 
}

type SaleHistory struct{

}

type Offer struct{
	Cost int
	Count int
}