package rabbit

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn     *amqp.Connection
	exchange string
	queue    string
	handler  MessageHandler
	mutex    sync.Mutex
}

type MessageHandler interface {
	HandleMessage(ctx context.Context, message []byte, routingKey string) error
}

type MessageHandlerFunc func(ctx context.Context, message []byte, routingKey string) error

func (f MessageHandlerFunc) HandleMessage(ctx context.Context, message []byte, routingKey string) error {
	return f(ctx, message, routingKey)
}

func NewConsumer(conn *amqp.Connection, exchange, queue string) *Consumer {
	return &Consumer{
		conn:     conn,
		exchange: exchange,
		queue:    queue,
	}
}

func (c *Consumer) SetHandler(handler MessageHandler) {
	c.handler = handler
}

func (c *Consumer) StartConsuming(ctx context.Context) error {
	if c.handler == nil {
		return errors.New("message handler not set")
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.conn.IsClosed() {
		return errors.New("connection is closed")
	}

	ch, err := c.conn.Channel()
	if err != nil {
		return fmt.Errorf("error creating channel: %w", err)
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		ch.Close()
		return fmt.Errorf("error setting QoS: %w", err)
	}

	msgs, err := ch.Consume(
		c.queue,
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		ch.Close()
		return fmt.Errorf("error starting consumer: %w", err)
	}

	go c.consumeMessages(ctx, ch, msgs)
	return nil
}

func (c *Consumer) consumeMessages(ctx context.Context, ch *amqp.Channel, msgs <-chan amqp.Delivery) {
	defer ch.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-msgs:
			if !ok {
				return
			}

			go c.processMessage(ctx, msg)
		}
	}
}

func (c *Consumer) processMessage(ctx context.Context, msg amqp.Delivery) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if c.handler == nil {
		msg.Nack(false, true)
		return
	}

	err := c.handler.HandleMessage(ctx, msg.Body, msg.RoutingKey)
	if err != nil {
		msg.Nack(false, true)
	} else {
		msg.Ack(false)
	}
}
