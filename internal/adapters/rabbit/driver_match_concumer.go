package rabbit

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"ride-hail/internal/core/domain/models"
	"ride-hail/pkg/rabbit"
)

type DriverResponseConsumer struct {
	consumer *rabbit.Consumer
	ch       chan models.DriverResponseEvent
}

const (
	driverResponseExchange = "driver_topic"
	driverResponseQueue    = "driver_responses"
)

func NewDriverResponseConsumer(conn *amqp.Connection) *DriverResponseConsumer {
	ch := make(chan models.DriverResponseEvent, 100)

	c := rabbit.NewConsumer(conn, driverResponseExchange, driverResponseQueue)

	c.SetHandler(rabbit.MessageHandlerFunc(func(ctx context.Context, msg []byte, rk string) error {
		var response models.DriverResponseEvent
		if err := json.Unmarshal(msg, &response); err != nil {
			return nil
		}

		select {
		case ch <- response:
		default:
		}
		return nil
	}))

	return &DriverResponseConsumer{
		consumer: c,
		ch:       ch,
	}
}

func (r *DriverResponseConsumer) Start(ctx context.Context) error {
	return r.consumer.StartConsuming(ctx)
}

func (r *DriverResponseConsumer) Subscribe(ctx context.Context) (<-chan models.DriverResponseEvent, error) {
	return r.ch, nil
}
