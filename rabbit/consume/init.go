package rabbit

import (
	"fmt"
	"time"
)

func Init() {
	go func() {
		for {
			err := ConsumeIncrementEvent()
			if err != nil {
				fmt.Errorf("%s", err.Error())
			}
			// logger.Info("RabbitMQ consumePlaceOrder conectando en 5 segundos.")
			time.Sleep(5 * time.Second)
		}
	}()

	go func() {
		for {
			err := ConsumeOrderPlaced()
			if err != nil {
				fmt.Errorf("%s", err.Error())
			}
			// logger.Info("RabbitMQ consumePlaceOrder conectando en 5 segundos.")
			time.Sleep(5 * time.Second)
		}
	}()
}
