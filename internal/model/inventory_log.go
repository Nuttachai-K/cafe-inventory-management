package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type InventoryLog struct {
	ID             int       `json:"id"`
	InventoryId    int       `json:"inventory_id"`
	UserId         int       `json:"user_id"`
	ChangeQuantity int       `json:"change_quantity"`
	Operation      string    `json:"operation"`
	CreatedAt      time.Time `json:"created_at"`
}

type InventoryLogFilter struct {
	InventoryId *int
	UserId      *int
	Operation   *string
	From, To    *time.Time
	Limit       int
}

type InventoryUpdateRequest struct {
	Operation      string `json:"operation"`
	ChangeQuantity int    `json:"change_quantity"`
}

type Operation string

const (
	InOperation     Operation = "IN"
	OutOperation    Operation = "OUT"
	AdjustOperation Operation = "ADJUST"
)

func (o Operation) Valid() bool {
	switch o {
	case InOperation, OutOperation, AdjustOperation:
		return true
	default:
		return false
	}
}

func ParseOperation(s string) (Operation, error) {
	o := Operation(strings.ToUpper(strings.TrimSpace(s)))
	if !o.Valid() {
		return "", fmt.Errorf("operation must be IN, OUT, or ADJUST: %q", s)
	}
	return o, nil
}

func (c InventoryLog) MarshalJSON() ([]byte, error) {
	type Alias InventoryLog
	return json.Marshal(&struct {
		*Alias
		CreatedAt string `json:"created_at"`
	}{
		Alias:     (*Alias)(&c),
		CreatedAt: formatDateHourJST(c.CreatedAt),
	})
}
