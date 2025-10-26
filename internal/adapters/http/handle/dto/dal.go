package dto

import (
	"errors"
	"fmt"

	"ride-hail/internal/core/domain/models"
	"ride-hail/internal/core/domain/types"
)

func ValidateDriver(d models.Driver) error {
	if d.ID == "" {
		return errors.New("driver ID (uuid string) is required")
	}
	if d.LicenseNumber == "" {
		return errors.New("license_number is required")
	}
	if len(d.LicenseNumber) > 50 {
		return errors.New("license_number too long (max 50)")
	}
	if d.Rating < 1.0 || d.Rating > 5.0 {
		return fmt.Errorf("rating must be between 1.0 and 5.0, got %v", d.Rating)
	}
	if d.TotalRides < 0 {
		return fmt.Errorf("total_rides must be >= 0, got %d", d.TotalRides)
	}
	if d.TotalEarnings < 0 {
		return fmt.Errorf("total_earnings must be >= 0, got %v", d.TotalEarnings)
	}
	if d.Status != types.DriverStatusAvailable || d.Status != types.DriverStatusBusy ||
		d.Status != types.DriverStatusEnRoute || d.Status != types.DriverStatusOffline {
		return fmt.Errorf("invalid status: %q", d.Status)
	}
	return nil
}

func ValidateDriverSession(s *models.DriverSession) error {
	if s.ID == "" {
		return errors.New("driver_session id is required")
	}
	if s.DriverID == "" {
		return errors.New("driver_id is required")
	}
	if s.TotalRides < 0 {
		return fmt.Errorf("total_rides must be >= 0, got %d", s.TotalRides)
	}
	if s.TotalEarnings < 0 {
		return fmt.Errorf("total_earnings must be >= 0, got %v", s.TotalEarnings)
	}
	return nil
}

func ValidateLocationHistory(lh *models.LocationHistory) error {
	if lh.ID == "" {
		return errors.New("location_history id is required")
	}
	if lh.Latitude < -90 || lh.Latitude > 90 {
		return fmt.Errorf("latitude must be between -90 and 90, got %v", lh.Latitude)
	}
	if lh.Longitude < -180 || lh.Longitude > 180 {
		return fmt.Errorf("longitude must be between -180 and 180, got %v", lh.Longitude)
	}
	if lh.HeadingDegrees != nil {
		if *lh.HeadingDegrees < 0 || *lh.HeadingDegrees > 360 {
			return fmt.Errorf("heading_degrees must be between 0 and 360, got %v", *lh.HeadingDegrees)
		}
	}
	return nil
}
