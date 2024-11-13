package requirereposition

import (
	"context"
	rabbitEmitter "emmanuel-guerreiro/stockgo/rabbit/emit"
	"encoding/json"
	"fmt"
)

func EmitNotEnoughStock(articleId string) error {
	ch, err := rabbitEmitter.GetChannel(context.Background())
	if err != nil {
		fmt.Println("Error getting channel require_reposition")
		return nil
	}

	if err = ch.ExchangeDeclare("require_reposition", "fanout"); err != nil {
		fmt.Println("Error declaring exchange require_reposition")
		return err
	}

	send := requireRepositionDto{
		ArticleId: articleId,
	}

	body, err := json.Marshal(send)
	if err != nil {
		fmt.Println("Error marshaling message require_reposition")
		return err
	}

	err = ch.Publish(
		"require_reposition", // exchange
		"require_reposition", // routing key
		body,
	)
	if err != nil {
		fmt.Println("Error publishing message require_reposition")

		return err
	}

	fmt.Println("Emited require_reposition")

	return nil
}

type requireRepositionDto struct {
	ArticleId string `bson:"articleId" json:"articleId" validate:"required,min=1,max=100"`
}
