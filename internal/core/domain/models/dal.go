package models

import (
	"time"
)

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

type Driver struct {
	ID            string       `json:"id"`
	CreatedAt     time.Time    `json:"creation_time"`
	UpdatedAt     time.Time    `json:"update_time"`
	LicenseNumber string       `json:"license_number"`
	VehicleType   string       `json:"vehicle_type"`
	VehicleAttrs  VehicleAttrs `json:"vehicle_attrs"`
	Rating        float64      `json:"rating"`
	TotalRides    int          `json:"total_rides"`
	TotalEarnings float64      `json:"total_earnings"`
	Status        string       `json:"status"`
	IsVarified    bool         `json:"is_varified"`
}
type VehicleAttrs struct {
	LicensePlate      string `json:"license_plate"`
	InspectionDate    string `json:"inspection_date"`
	Make              string `json:"make"`
	Model             string `json:"model"`
	Year              int    `json:"year"`
	Color             string `json:"color"`
	Seats             int    `json:"seats"`
	InsuranceExpiry   string `json:"insurance_expiry"`
	TaxiLicenseExpiry string `json:"taxi_license_expiry"`
}

type DriverSession struct {
	ID            string    `json:"id"`
	DriverID      string    `json:"driver_id"`
	StartedAt     time.Time `json:"started_at"`
	EndedAt       time.Time `json:"ended_at"`
	TotalRides    int       `json:"total_rides"`
	TotalEarnings float64   `json:"total_earnings"`
}

type SessionSummary struct {
	DurationHours  float64 `json:"duration_hours"`
	RidesCompleted int     `json:"rides_completed"`
	Earnings       float64 `json:"earnings"`
}

type DriverInfoClosed struct {
	Status         string         `json:"status"`
	SessionID      string         `json:"session_id"`
	SessionSummary SessionSummary `json:"session_summary"`
	Message        string         `json:"message"`
}
