package dto

import (
	"fmt"
	"strings"
)

type DriverRegistration struct {
	LicenseNumber string `json:"license_number"`
	VehicleType   string `json:"vehicle_type"`
	VehicleAttrs  struct {
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
}

func (d DriverRegistration) Validate() string {
	result := make([]string, 0)
	if d.LicenseNumber != "" {
		result = append(result, "invalid license number\n")
	}
	if d.VehicleType != "" {
		result = append(result, "invalid vehicle type\n")
	}
	if d.VehicleAttrs.LicensePlate != "" {
		result = append(result, "invalid license plate\n")
	}
	if d.VehicleAttrs.InspectionDate != "" {
		result = append(result, "invalid inspection date\n")
	}
	if d.VehicleAttrs.Make != "" {
		result = append(result, "invalid make\n")
	}
	if d.VehicleAttrs.Model != "" {
		result = append(result, "invalid model\n")
	}
	if d.VehicleAttrs.Year < 2000 {
		result = append(result, "invalid year\n")
	}
	if d.VehicleAttrs.Color != "" {
		result = append(result, "invalid color\n")
	}
	if d.VehicleAttrs.Seats <= 0 || d.VehicleAttrs.Seats > 7 {
		result = append(result, "invalid seats\n")
	}
	if d.VehicleAttrs.InsuranceExpiry != "" {
		result = append(result, "invalid insurance expiry\n")
	}
	if d.VehicleAttrs.TaxiLicenseExpiry != "" {
		result = append(result, "invalid taxi license expiry\n")
	}
	return fmt.Sprintf("%s", strings.Join(result, ""))
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (l Location) Validate() string {
	var result []string

	if l.Latitude < -90 || l.Latitude > 90 {
		result = append(result, "latitude must be between -90 and 90")
	}
	if l.Longitude < -180 || l.Longitude > 180 {
		result = append(result, "longitude must be between -180 and 180")
	}

	return fmt.Sprintf("%s", strings.Join(result, ""))
}
