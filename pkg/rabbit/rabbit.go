package rabbit

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Rabbit struct {
	conn *amqp.Connection
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
	conn, err := amqp.Dial(cfg.GetRabbitDsn())
	if err != nil {
		return nil, err
	}
	return &Rabbit{conn: conn, cfg: cfg}, nil
}

func (r *Rabbit) Close() {
	if r.conn != nil {
		_ = r.conn.Close()
	}
}

type QueueConfig struct {
	Name       string
	RoutingKey string
}

func (r *Rabbit) SetupExchangesAndQueues(exchangeName, exchangeType string, queues []QueueConfig) error {
	ch, err := r.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	if err = r.ensureExchange(ch, exchangeName, exchangeType); err != nil {
		return err
	}

	for _, qCfg := range queues {
		q, err := r.ensureQueue(ch, qCfg.Name)
		if err != nil {
			return err
		}

		if err = ch.QueueBind(
			q.Name,
			qCfg.RoutingKey,
			exchangeName,
			false,
			nil,
		); err != nil {
			return err
		}
	}

	return nil
}

func (r *Rabbit) ensureExchange(ch *amqp.Channel, name, kind string) error {
	if err := ch.ExchangeDeclarePassive(name, kind, true, false, false, false, nil); err == nil {
		return nil
	}

	return ch.ExchangeDeclare(name, kind, true, false, false, false, nil)
}

func (r *Rabbit) ensureQueue(ch *amqp.Channel, name string) (amqp.Queue, error) {
	_, err := ch.QueueDeclarePassive(name, true, false, false, false, nil)
	if err == nil {
		return amqp.Queue{Name: name}, nil
	}

	q, err := ch.QueueDeclare(name, true, false, false, false, nil)
	if err != nil {
		return amqp.Queue{}, err
	}
	return q, nil
}
