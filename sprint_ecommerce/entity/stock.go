package entity

type Stock struct {
	ID			string		`json:"id"`
	Stock 		int			`json:"stock"`
}

type StockUpdateRequest struct {
	ProductID	*string		`json:"product_id"`
	Stock 		*int			`json:"stock"`
}