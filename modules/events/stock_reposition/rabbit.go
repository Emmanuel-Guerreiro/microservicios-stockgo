package stockreposition

import (
	"emmanuel-guerreiro/stockgo/lib"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ConsumeRepositionEvent() error {
	conn, err := amqp.Dial(lib.GetEnv().RabbitURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"stock_reposition", // name
		"direct",           // type
		false,              // durable
		false,              // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		fmt.Errorf("%s", err.Error())
		return err
	}

	q, err := ch.QueueDeclare(
		"stock_reposition", // name
		false,              // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		return err
	}

	err = ch.QueueBind(
		q.Name,             // queue name
		"stock_reposition", // routing key
		"stock_reposition", // exchange
		false,
		nil)
	if err != nil {
		// logger.Error(err)
		return err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack //TODo: Should be handled manually?
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	go func() {
		for d := range msgs {
			newMessage := &consumeStockRepositionDto{}
			body := d.Body

			err = json.Unmarshal(body, newMessage)
			if err != nil {
				fmt.Errorf("%s", err.Error())
				//TODO: Should requeue message?
				continue
			}

			handleReposition(newMessage)
		}
	}()

	fmt.Println("Closed connection: ", <-conn.NotifyClose(make(chan *amqp.Error)))
	return nil
}

func ListenerReposition() {
	for {
		err := ConsumeRepositionEvent()
		if err != nil {
			fmt.Errorf("%s", err.Error())
		}
		// logger.Info("RabbitMQ consumePlaceOrder conectando en 5 segundos.")
		time.Sleep(5 * time.Second)
	}
}
