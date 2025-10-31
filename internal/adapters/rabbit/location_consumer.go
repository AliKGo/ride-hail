package rabbit

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"ride-hail/internal/core/domain/models"
	"ride-hail/pkg/rabbit"
)

type LocationConsumer struct {
	consumer *rabbit.Consumer
	ch       chan models.DriverLocationUpdate
}

const (
	exName    = "location_fanout"
	queueName = "location_updates_ride"
)

func NewLocationConsumer(conn *amqp.Connection) *LocationConsumer {
	ch := make(chan models.DriverLocationUpdate, 100)
	c := rabbit.NewConsumer(conn, exName, queueName)
	c.SetHandler(rabbit.MessageHandlerFunc(func(ctx context.Context, msg []byte, rk string) error {
		var loc models.DriverLocationUpdate
		if err := json.Unmarshal(msg, &loc); err != nil {
			return nil
		}
		select {
		case ch <- loc:
		default:
		}
		return nil
	}))
	return &LocationConsumer{consumer: c, ch: ch}
}

func (r *LocationConsumer) Start(ctx context.Context) error {
	return r.consumer.StartConsuming(ctx)
}

func (r *LocationConsumer) Subscribe(ctx context.Context) (<-chan models.DriverLocationUpdate, error) {
	return r.ch, nil
}
