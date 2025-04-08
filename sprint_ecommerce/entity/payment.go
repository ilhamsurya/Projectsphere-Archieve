package entity

type Payment struct {
	PaymentProofImgUrl string `json:"paymentProofImageUrl"`
	ID                 int   `json:"id"`
	ProductID          string   `json:"productId"`
	BankAccountID      string   `json:"bankAccountId"`
	Quantity           int   `json:"quantity"`
}

type PaymentCreateRequest struct {
	PaymentProofImgUrl *string `json:"paymentProofImageUrl"`
	ProductID          *string   `json:"productId"`
	BankAccountID      *string   `json:"bankAccountId"`
	Quantity           *int   `json:"quantity"`
}
