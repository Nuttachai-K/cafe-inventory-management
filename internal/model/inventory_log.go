package model

import "time"

type InventoryLog struct {
	ID             int       `json:"id"`
	InventoryId    int       `json:"inventory_id"`
	UserId         int       `json:"user_id"`
	ChangeQuantity int       `json:"change_quantity"`
	Operation      string    `json:"operation"`
	CreatedAt      time.Time `json:"created_at"`
}
