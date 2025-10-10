package rabbit

import (
	"fmt"
	"github.com/rabbitmq/amqp091-go"
)

type Rabbit struct {
	conn *amqp091.Connection
	cfg  Config
}

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
}

func (c Config) GetRabbitDsn() string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%d/",
		c.User,
		c.Password,
		c.Host,
		c.Port,
	)
}

func New(cfg Config) (*Rabbit, error) {
	conn, err := amqp091.Dial(cfg.GetRabbitDsn())
	if err != nil {
		return nil, err
	}
	return &Rabbit{conn: conn, cfg: cfg}, nil
}
