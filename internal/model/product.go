package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Product struct {
	ID          int             `json:"id"`
	CafeId      int             `json:"cafe_id"`
	CategoryId  int             `json:"category_id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Price       decimal.Decimal `json:"price"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
type ProductUpdate struct {
	CafeId      *int             `json:"cafe_id"`
	CategoryId  *int             `json:"category_id"`
	Name        *string          `json:"name"`
	Description *string          `json:"description"`
	Price       *decimal.Decimal `json:"price"`
}

type ProductWithCategory struct {
	Product
	CategoryName string `json:"category_name"`
}
