package rabbitEmitter

import (
	"emmanuel-guerreiro/stockgo/lib"
	"emmanuel-guerreiro/stockgo/lib/log"
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitChannel interface {
	ExchangeDeclare(name string, chType string) error
	Publish(exchange string, routingKey string, body []byte) error
}

type rabbitChannel struct {
	ch *amqp.Channel
}

// ErrChannelNotInitialized Rabbit channel could not be initialized
var ErrChannelNotInitialized = errors.New("channel not initialized")

func GetChannel(ctx ...interface{}) (RabbitChannel, error) {
	for _, o := range ctx {
		if ti, ok := o.(RabbitChannel); ok {
			return ti, nil
		}
	}

	conn, err := amqp.Dial(lib.GetEnv().RabbitURL)
	if err != nil {
		log.Get(ctx...).Error(err)
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Get(ctx...).Error(err)
		return nil, err
	}

	return rabbitChannel{ch: channel}, nil
}

func (c rabbitChannel) ExchangeDeclare(
	name string,
	chType string,
) error {
	return c.ch.ExchangeDeclare(
		name,   // name
		chType, // type
		false,  // durable
		false,  // auto-deleted
		false,  // internal
		false,  // no-wait
		nil,    // arguments
	)
}
func (c rabbitChannel) Publish(
	exchange string,
	routingKey string,
	body []byte,
) error {
	return c.ch.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			Body: body,
		})
}
