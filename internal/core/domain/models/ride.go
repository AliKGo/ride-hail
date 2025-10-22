package models

import "time"

type CreateRideRequest struct {
	PassengerID          string
	PickupLatitude       float64
	PickupLongitude      float64
	PickupAddress        string
	DestinationLatitude  float64
	DestinationLongitude float64
	DestinationAddress   string
	RideType             string
}

type Ride struct {
	ID                      string    `json:"id"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
	RideNumber              string    `json:"ride_number"`
	PassengerID             string    `json:"passenger_id"`
	DriverID                string    `json:"driver_id"`
	VehicleType             string    `json:"vehicle_type"`
	Status                  string    `json:"status"`
	Priority                int       `json:"priority"`
	RequestedAt             time.Time `json:"requested_at"`
	MatchedAt               time.Time `json:"matched_at"`
	ArrivedAt               time.Time `json:"arrival_at"`
	StartedAt               time.Time `json:"started_at"`
	CompletedAt             time.Time `json:"completed_at"`
	CancelledAt             time.Time `json:"cancelled_at"`
	CancellationReason      string    `json:"cancellation_reason"`
	EstimatedFare           float64   `json:"estimated_fare"`
	FinalFare               float64   `json:"final_fare"`
	PickupCoordinateId      string    `json:"pickup_coordinate_id"`
	DestinationCoordinateId string    `json:"destination_coordinate_id"`
}

type CreateRideResponse struct {
	RideID                   string  `json:"ride_id"`
	RideNumber               string  `json:"ride_number"`
	Status                   string  `json:"status"`
	EstimatedFare            float64 `json:"estimated_fare"`
	EstimatedDurationMinutes int     `json:"estimated_duration_minutes"`
	EstimatedDistanceKm      float64 `json:"estimated_distance_km"`
}

type CloseRideRequest struct {
	RideID string `json:"ride_id"`
	Reason string `json:"reason"`
}

type CloseRideResponse struct {
	RideID      string    `json:"ride_id"`
	Status      string    `json:"status"`
	CancelledAt time.Time `json:"cancelled_at"`
	Message     string    `json:"message"`
}
