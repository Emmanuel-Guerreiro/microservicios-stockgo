package stockreposition

import (
	"context"
	"emmanuel-guerreiro/stockgo/lib"
	rabbit "emmanuel-guerreiro/stockgo/rabbit/emit"
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
		fmt.Println(err.Error())
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
				fmt.Println(err.Error())
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
			fmt.Println(err.Error())
		}
		// logger.Info("RabbitMQ consumePlaceOrder conectando en 5 segundos.")
		time.Sleep(5 * time.Second)
	}
}

func emitStockNowAvailable(articleId string, correlationId string) error {
	ch, err := rabbit.GetChannel(context.Background())
	if err != nil {
		fmt.Println("Error getting channel stock_available")
		return nil
	}

	if err = ch.ExchangeDeclare("stock_available", "fanout"); err != nil {
		fmt.Println("Error declaring exchange stock_available")
		return err
	}

	send := placeStockAvailableMessageDto{
		ArticleId:     articleId,
		CorrelationId: correlationId,
	}

	body, err := json.Marshal(send)
	if err != nil {
		fmt.Println("Error marshaling message stock_available")
		return err
	}

	err = ch.Publish(
		"stock_available", // exchange
		"",                // routing key -> No lleva RK porque es fanout el exchange al que me estoy suscribiendo
		body,
	)
	if err != nil {
		fmt.Println("Error publishing message stock_available")

		return err
	}

	fmt.Println("Emited stock_available", string(body))

	return nil
}
