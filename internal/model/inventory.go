package model

import (
	"encoding/json"
	"time"
)

type Inventory struct {
	ID            int       `json:"id"`
	ProductId     int       `json:"product_id"`
	StockQuantity int       `json:"stock_quantity"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type InventoryWithProduct struct {
	ID            int       `json:"id"`
	ProductId     int       `json:"product_id"`
	ProductName   string    `json:"product_name"`
	CafeId        int       `json:"cafe_id"`
	CafeName      string    `json:"cafe_name"`
	StockQuantity int       `json:"stock_quantity"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (i Inventory) MarshalJSON() ([]byte, error) {
	type Alias Inventory
	return json.Marshal(&struct {
		*Alias
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}{
		Alias:     (*Alias)(&i),
		CreatedAt: formatDateHourJST(i.CreatedAt),
		UpdatedAt: formatDateHourJST(i.UpdatedAt),
	})
}

func (i InventoryWithProduct) MarshalJSON() ([]byte, error) {
	type Alias InventoryWithProduct
	return json.Marshal(&struct {
		*Alias
		UpdatedAt string `json:"updated_at"`
	}{
		Alias:     (*Alias)(&i),
		UpdatedAt: formatDateHourJST(i.UpdatedAt),
	})
}
