package model

import "time"

// Cafe represent a cafe in the system
type Cafe struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Address        string    `json:"address"`
	Latitude       float64   `json:"latitude"`
	Longitude      float64   `json:"longtitude"`
	NearestStation string    `json:"nearest_station"`
	OpeningTime    time.Time `json:"opening_time"`
	ClosingTime    time.Time `json:"closing_time"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
