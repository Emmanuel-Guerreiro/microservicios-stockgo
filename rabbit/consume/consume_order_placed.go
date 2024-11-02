package rabbit

import (
	"emmanuel-guerreiro/stockgo/lib"
	ordersplaced "emmanuel-guerreiro/stockgo/modules/orders_placed"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ConsumeOrderPlaced() error {

	conn, err := amqp.Dial(lib.GetEnv().RabbitURL)
	if err != nil {
		return err
	}
	defer conn.Close()
	chn, err := conn.Channel()
	if err != nil {
		return err
	}
	defer chn.Close()

	err = chn.ExchangeDeclare(
		"order_placed", // name
		"fanout",       // type
		false,          // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return err
	}

	queue, err := chn.QueueDeclare(
		"catalog_order_placed", // name
		false,                  // durable
		false,                  // delete when unused
		false,                  // exclusive
		false,                  // no-wait
		nil,                    // arguments
	)
	if err != nil {
		// logger.Error(err)
		return err
	}

	err = chn.QueueBind(
		queue.Name,     // queue name
		"",             // routing key
		"order_placed", // exchange
		false,
		nil)
	if err != nil {
		// logger.Error(err)
		return err
	}

	mgs, err := chn.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		// logger.Error(err)
		return err
	}

	fmt.Println("RabbitMQ consumeOrderPlaced conectado")
	go func() {
		for d := range mgs {
			body := d.Body

			articleMessage := &ordersplaced.ConsumeOrderPlacedDto{}
			err = json.Unmarshal(body, articleMessage)
			if err == nil {
				// l := logger.WithField(log.LOG_FIELD_CORRELATION_ID, getConsumeOrderPlacedCorrelationId(articleMessage))
				// l.Info("Incoming order_placed :", string(body))

				ordersplaced.ProcessOrderPlaced(articleMessage)

				if err := d.Ack(false); err != nil {
					// l.Info("Failed ACK order_placed :", string(body), err)
					fmt.Println("Failed ACK order_placed :  ", string(body))
				} else {
					fmt.Println("Consumed order_placed :", string(body))

					// l.Info("Consumed order_placed :", string(body))
				}
			} else {
				fmt.Println("Error during parse")
				// logger.Error(err)
			}
		}
	}()

	fmt.Println("Closed connection: ", <-conn.NotifyClose(make(chan *amqp.Error)))

	return nil
}
