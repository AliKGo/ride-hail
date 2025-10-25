package models

import (
	"encoding/json"
	"time"
)

type Driver struct {
	ID            string          `db:"id"`             // uuid
	CreatedAt     time.Time       `db:"created_at"`     // timestamptz
	UpdatedAt     time.Time       `db:"updated_at"`     // timestamptz
	LicenseNumber string          `db:"license_number"` // varchar(50)
	VehicleType   *string         `db:"vehicle_type"`   // text
	VehicleAttrs  json.RawMessage `db:"vehicle_attrs"`  // jsonb
	Rating        float64         `db:"rating"`         // decimal(3,2)
	TotalRides    int             `db:"total_rides"`    // integer
	TotalEarnings float64         `db:"total_earnings"` // decimal(10,2)
	Status        string          `db:"status"`         // text -> driver_status
	IsVerified    bool            `db:"is_verified"`    // boolean
}

type DriverSession struct {
	ID            string     `db:"id"`             // uuid
	DriverID      string     `db:"driver_id"`      // uuid
	StartedAt     time.Time  `db:"started_at"`     // timestamptz
	EndedAt       *time.Time `db:"ended_at"`       // nullable timestamptz
	TotalRides    int        `db:"total_rides"`    // integer
	TotalEarnings float64    `db:"total_earnings"` // decimal(10,2)
}

type LocationHistory struct {
	ID             string    `db:"id"`              // uuid
	CoordinateID   *string   `db:"coordinate_id"`   // uuid nullable
	DriverID       *string   `db:"driver_id"`       // uuid nullable
	Latitude       float64   `db:"latitude"`        // decimal(10,8)
	Longitude      float64   `db:"longitude"`       // decimal(11,8)
	AccuracyMeters *float64  `db:"accuracy_meters"` // decimal(6,2)
	SpeedKmh       *float64  `db:"speed_kmh"`       // decimal(5,2)
	HeadingDegrees *float64  `db:"heading_degrees"` // decimal(5,2)
	RecordedAt     time.Time `db:"recorded_at"`     // timestamptz
	RideID         *string   `db:"ride_id"`         // uuid nullable
}
