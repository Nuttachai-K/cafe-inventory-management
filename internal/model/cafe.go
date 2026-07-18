package model

import "time"

// Cafe represent a cafe in the system
type Cafe struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Address        string    `json:"address"`
	Latitude       float64   `json:"latitude"`
	Longitude      float64   `json:"longitude"`
	NearestStation string    `json:"nearest_station"`
	OpeningTime    time.Time `json:"opening_time"`
	ClosingTime    time.Time `json:"closing_time"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CafeUpdate struct {
	Name           *string    `json:"name,omitempty"`
	Address        *string    `json:"address,omitempty"`
	Latitude       *float64   `json:"latitude,omitempty"`
	Longitude      *float64   `json:"longitude,omitempty"`
	NearestStation *string    `json:"nearest_station,omitempty"`
	OpeningTime    *time.Time `json:"opening_time,omitempty"`
	ClosingTime    *time.Time `json:"closing_time,omitempty"`
}
