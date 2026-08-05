package model

import "time"

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
