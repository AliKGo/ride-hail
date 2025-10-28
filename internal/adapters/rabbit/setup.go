package rabbit

import (
	"errors"
	"ride-hail/pkg/rabbit"
)

func InitRabbitTopology(r *rabbit.Rabbit) error {
	if r.Conn.IsClosed() {
		return errors.New("connection is closed")
	}

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

	if err := r.SetupExchangesAndQueues(exchanges[0].Name, exchanges[0].Type, rideQueues); err != nil {
		return err
	}
	if err := r.SetupExchangesAndQueues(exchanges[1].Name, exchanges[1].Type, driverQueues); err != nil {
		return err
	}
	if err := r.SetupExchangesAndQueues(exchanges[2].Name, exchanges[2].Type, locationQueues); err != nil {
		return err
	}

	return nil
}
