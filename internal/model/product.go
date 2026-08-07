package model

import (
	"encoding/json"
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
	ID           int             `json:"id"`
	CafeId       int             `json:"cafe_id"`
	CategoryId   int             `json:"category_id"`
	CategoryName string          `json:"category_name"`
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	Price        decimal.Decimal `json:"price"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

func (p Product) MarshalJSON() ([]byte, error) {
	type Alias Product
	return json.Marshal(&struct {
		*Alias
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}{
		Alias:     (*Alias)(&p),
		CreatedAt: formatDateHourJST(p.CreatedAt),
		UpdatedAt: formatDateHourJST(p.UpdatedAt),
	})
}

func (p ProductWithCategory) MarshalJSON() ([]byte, error) {
	type Alias ProductWithCategory
	return json.Marshal(&struct {
		*Alias
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}{
		Alias:     (*Alias)(&p),
		CreatedAt: formatDateHourJST(p.CreatedAt),
		UpdatedAt: formatDateHourJST(p.UpdatedAt),
	})
}
