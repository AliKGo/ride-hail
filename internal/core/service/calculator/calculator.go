package calculator

import (
	"fmt"
	"math"
	"ride-hail/internal/core/domain/types"
)

const earthRadius = 6371.0

const avgSpeedKmH = 30 // km/h

func Distance(lat1, lon1, lat2, lon2 float64) float64 {
	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLon := (lon2 - lon1) * math.Pi / 180.0

	lat1Rad := lat1 * math.Pi / 180.0
	lat2Rad := lat2 * math.Pi / 180.0

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1Rad)*math.Cos(lat2Rad)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

func Duration(dist float64) int {
	return int((dist / avgSpeedKmH) * 60)
}

func CalculateFare(rideType string, distanceKm float64, durationMin int) (float64, error) {
	var baseFare, ratePerKm, ratePerMin float64

	switch rideType {
	case types.RideTypeECONOMY:
		baseFare = 500
		ratePerKm = 100
		ratePerMin = 50
	case types.RideTypePREMIUM:
		baseFare = 800
		ratePerKm = 120
		ratePerMin = 60
	case types.RideTypeXL:
		baseFare = 1000
		ratePerKm = 150
		ratePerMin = 75
	default:
		return 0, fmt.Errorf("unknown ride type: %s", rideType)
	}

	total := baseFare + (distanceKm * ratePerKm) + (float64(durationMin) * ratePerMin)
	return total, nil
}
