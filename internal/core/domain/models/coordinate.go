package models

import "time"

type Coordinate struct {
	ID              string    `json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	EntityID        string    `json:"entity_id"`
	EntityType      string    `json:"entity_type"`
	Address         string    `json:"address"`
	Latitude        float64   `json:"latitude"`
	Longitude       float64   `json:"longitude"`
	FareAmount      float64   `json:"fare_amount"`
	DistanceKM      float64   `json:"distance_km"`
	DurationMinutes int       `json:"duration_minutes"`
	IsCurrent       bool      `json:"is_current"`
}
