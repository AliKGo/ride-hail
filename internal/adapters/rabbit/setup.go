package rabbit

import (
	"ride-hail/pkg/rabbit"
)

func InitRabbitTopology(r *rabbit.Rabbit) error {
	exchanges := []struct {
		Name string
		Type string
	}{
		{"ride_topic", "topic"},
		{"driver_topic", "topic"},
		{"location_fanout", "fanout"},
	}

	for _, ex := range exchanges {
		if err := r.SetupExchangesAndQueues(ex.Name, ex.Type, nil); err != nil {
			return err
		}
	}

	rideQueues := []rabbit.QueueConfig{
		{"ride_requests", "ride.request.*"},
		{"ride_status", "ride.status.*"},
	}
	driverQueues := []rabbit.QueueConfig{
		{"driver_matching", "driver.request.*"},
		{"driver_responses", "driver.response.*"},
		{"driver_status", "driver.status.*"},
	}
	locationQueues := []rabbit.QueueConfig{
		{"location_updates_ride", ""},
	}

	if err := r.SetupExchangesAndQueues("ride_topic", "topic", rideQueues); err != nil {
		return err
	}
	if err := r.SetupExchangesAndQueues("driver_topic", "topic", driverQueues); err != nil {
		return err
	}
	if err := r.SetupExchangesAndQueues("location_fanout", "fanout", locationQueues); err != nil {
		return err
	}

	return nil
}
