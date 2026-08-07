package model

import (
	"encoding/json"
	"time"
)

type Category struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c Category) MarshalJSON() ([]byte, error) {
	type Alias Category
	return json.Marshal(&struct {
		*Alias
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}{
		Alias:     (*Alias)(&c),
		CreatedAt: formatDateHourJST(c.CreatedAt),
		UpdatedAt: formatDateHourJST(c.UpdatedAt),
	})
}
