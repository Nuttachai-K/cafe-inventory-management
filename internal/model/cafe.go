package model

import (
	"encoding/json"
	"time"
)

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
	DistanceKm     *float64  `json:"distance_km,omitempty"`
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

// For swagger documentation
type CafeCreate struct {
	Name           string    `json:"name"`
	Address        string    `json:"address"`
	Latitude       float64   `json:"latitude"`
	Longitude      float64   `json:"longitude"`
	NearestStation string    `json:"nearest_station"`
	OpeningTime    time.Time `json:"opening_time"`
	ClosingTime    time.Time `json:"closing_time"`
}

type CafeFilter struct {
	Station  *string
	Limit    int
	Lat      *float64
	Lng      *float64
	RadiusKm *float64
}

func (c Cafe) MarshalJSON() ([]byte, error) {
	type Alias Cafe
	return json.Marshal(&struct {
		*Alias
		OpeningTime string `json:"opening_time"`
		ClosingTime string `json:"closing_time"`
		CreatedAt   string `json:"created_at"`
		UpdatedAt   string `json:"updated_at"`
	}{
		Alias:       (*Alias)(&c),
		OpeningTime: formatHour(c.OpeningTime),
		ClosingTime: formatHour(c.ClosingTime),
		CreatedAt:   formatDateHourJST(c.CreatedAt),
		UpdatedAt:   formatDateHourJST(c.UpdatedAt),
	})
}

func (c *Cafe) UnmarshalJSON(data []byte) error {
	type Alias Cafe
	aux := &struct {
		OpeningTime string `json:"opening_time"`
		ClosingTime string `json:"closing_time"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	if aux.OpeningTime != "" {
		t, err := parseClockTime(aux.OpeningTime)
		if err != nil {
			return err
		}
		c.OpeningTime = t
	}
	if aux.ClosingTime != "" {
		t, err := parseClockTime(aux.ClosingTime)
		if err != nil {
			return err
		}
		c.ClosingTime = t
	}
	return nil
}

func (cu *CafeUpdate) UnmarshalJSON(data []byte) error {
	type Alias CafeUpdate
	aux := &struct {
		OpeningTime string `json:"opening_time"`
		ClosingTime string `json:"closing_time"`
		*Alias
	}{
		Alias: (*Alias)(cu),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	if aux.OpeningTime != "" {
		t, err := parseClockTime(aux.OpeningTime)
		if err != nil {
			return err
		}
		cu.OpeningTime = &t
	}
	if aux.ClosingTime != "" {
		t, err := parseClockTime(aux.ClosingTime)
		if err != nil {
			return err
		}
		cu.ClosingTime = &t
	}
	return nil
}
