package entity

import (
	"github.com/google/uuid"
)

type Product struct {
	ID            uuid.UUID `db:"id" json:"productId"`
	UserID        uuid.UUID `db:"user_id" json:"-"`
	Name          string    `db:"name" json:"name" validate:"required,max=60,min=5"`
	Price         float64   `db:"price" json:"price" validate:"required,numeric,gte=0"`
	Stock         int       `db:"stock" json:"stock" validate:"required,numeric,gte=0"`
	ImageUrl      string    `db:"image_url" json:"imageUrl" validate:"required,url"`
	Condition     string    `db:"condition" json:"condition" validate:"required,validCondition"`
	IsPurchasable bool      `db:"is_purchasable" json:"isPurchasable" validate:"required"`
	Tags          []string  `db:"tags" json:"tags" validate:"required,min=0"`
	PurchaseCount int       `json:"purchaseCount"`
}

type Key struct {
	UserOnly       *bool    `json:"userOnly"`
	Limit          *int     `json:"limit"`
	Offset         *int     `json:"offset"`
	Tags           []string `json:"tags"`
	Condition      *string  `json:"condition"`
	ShowEmptyStock *bool    `json:"showEmptyStock"`
	MaxPrice       *float64 `json:"maxPrice"`
	MinPrice       *float64 `json:"minPrice"`
	SortBy         *string  `json:"sortBy"`
	OrderBy        *string  `json:"orderBy"`
	Search         *string  `json:"search"`
}

type ProductPayment struct {
	Name      string `json:"name"`
	TotalSold int    `json:"totalSold"`
}
