package rabbit

import (
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"
)

type Producer struct {
	conn  *amqp.Connection
	mutex sync.Mutex
}

func NewPublisher(conn *amqp.Connection) *Producer {
	return &Producer{
		conn: conn,
	}
}

func (p *Producer) Producer(exName, queue string, message []byte) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.conn.IsClosed() {
		return errors.New("connection is closed")
	}

	ch, err := p.conn.Channel()
	if err != nil {
		return fmt.Errorf("error in creating channel %w", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf("error in declaring queue %w", err)
	}

	err = ch.Publish(
		exName,
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)

	if err != nil {
		return fmt.Errorf("error in publishing message %w", err)
	}
	return nil
}
