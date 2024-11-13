package stockviews

import (
	"emmanuel-guerreiro/stockgo/lib"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ConsumeStockConsultEvent() error {
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
		"stock_consulting", // name
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
		"get_stock", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		return err
	}

	err = ch.QueueBind(
		q.Name,             // queue name
		"get_stock",        // routing key
		"stock_consulting", // exchange
		false,
		nil,
	)

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
			fmt.Printf("Received a message: %s\n", d.Body)
			// fmt.Printf("Message: %s\n", d.Body)
			// fmt.Printf("Message properties: %v\n", d.Properties)
			// fmt.Printf("Delivery tag: %d\n", d.DeliveryTag)
			// fmt.Printf("Ack: %v\n", d.Ack)
			// fmt.Printf("Nack: %v\n", d.Nack)
		}
	}()

	fmt.Println("Closed connection: ", <-conn.NotifyClose(make(chan *amqp.Error)))
	return nil
}

func ListenerStockConsultEvent() {
	for {
		err := ConsumeStockConsultEvent()
		if err != nil {
			fmt.Errorf("%s", err.Error())
		}
		// logger.Info("RabbitMQ consumePlaceOrder conectando en 5 segundos.")
		time.Sleep(5 * time.Second)
	}
}
