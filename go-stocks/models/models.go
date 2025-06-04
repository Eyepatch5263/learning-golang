package models

type Stock struct {
	StockID  int64     `json:"stockid"`
	StockName string    `json:"name"`
	StockPrice float64   `json:"price"`
	StockCompany string    `json:"company"`
}



