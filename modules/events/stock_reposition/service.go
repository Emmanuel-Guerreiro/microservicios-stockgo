package stockreposition

import (
	"emmanuel-guerreiro/stockgo/modules/events"
	stockviews "emmanuel-guerreiro/stockgo/modules/stock_views"
	"fmt"
)

func handleReposition(message *consumeStockRepositionDto) {

	incDto := &events.CreateEventDto{
		Type: events.Reposition,
		RepositionEvent: &events.RepositionEvent{
			ArticleId: message.Message.ArticleId,
			Quantity:  message.Message.Amount,
		},
	}

	if _, err := events.CreateEvent(incDto); err != nil {
		//TODO: Place an event in the event bus notifying that the order cant be processed
		fmt.Println("ERROR AL CREAR EVENTO", err)
		return
	}

	if _, err := stockviews.GenerateStockView(message.Message.ArticleId); err != nil {
		//TODO:Should notify to somewhere?
		fmt.Println("ERROR AL REGENERAR STOCKVIEWS", message.Message.ArticleId)
	}

	return
}
