package models

import "time"

type DriverLocationUpdate struct {
	DriverID  string         `json:"driver_id"`
	RideID    string         `json:"ride_id"`
	Location  LocationDriver `json:"location"`
	SpeedKmh  float64        `json:"speed_kmh"`
	Heading   float64        `json:"heading_degrees"`
	Timestamp time.Time      `json:"timestamp"`
}

type LocationDriver struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type DriverResponseEvent struct {
	RideID                  string `json:"ride_id"`
	DriverID                string `json:"driver_id"`
	Accepted                bool   `json:"accepted"`
	EstimatedArrivalMinutes int    `json:"estimated_arrival_minutes"`
	DriverLocation          struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	} `json:"driver_location"`
	DriverInfo struct {
		Name    string  `json:"name"`
		Rating  float64 `json:"rating"`
		Vehicle struct {
			Make  string `json:"make"`
			Model string `json:"model"`
			Color string `json:"color"`
			Plate string `json:"plate"`
		} `json:"vehicle"`
	} `json:"driver_info"`
	CorrelationID string `json:"correlation_id"`
}
