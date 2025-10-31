package rabbit

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"ride-hail/internal/core/domain/models"
	"ride-hail/pkg/rabbit"
)

type RideStatusConsumer struct {
	consumer *rabbit.Consumer
	ch       chan models.RideStatusEvent
}

const (
	rideStatusExchange = "ride_topic"
	rideStatusQueue    = "ride_statuses"
)

func NewRideStatusConsumer(conn *amqp.Connection) *RideStatusConsumer {
	ch := make(chan models.RideStatusEvent, 100)

	c := rabbit.NewConsumer(conn, rideStatusExchange, rideStatusQueue)

	c.SetHandler(rabbit.MessageHandlerFunc(func(ctx context.Context, msg []byte, rk string) error {
		var event models.RideStatusEvent
		if err := json.Unmarshal(msg, &event); err != nil {
			fmt.Printf("failed to unmarshal ride status event: %v\n", err)
			return nil
		}

		select {
		case ch <- event:
		default:
			fmt.Println("ride status channel full, dropping message")
		}
		return nil
	}))

	return &RideStatusConsumer{
		consumer: c,
		ch:       ch,
	}
}

func (r *RideStatusConsumer) Start(ctx context.Context) error {
	return r.consumer.StartConsuming(ctx)
}

func (r *RideStatusConsumer) Subscribe(ctx context.Context) (<-chan models.RideStatusEvent, error) {
	return r.ch, nil
}
