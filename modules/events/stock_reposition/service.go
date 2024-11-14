package stockreposition

import (
	"emmanuel-guerreiro/stockgo/modules/events"
	stockviews "emmanuel-guerreiro/stockgo/modules/stock_views"
	"fmt"
)

func handleReposition(message *consumeStockRepositionDto) {
	fmt.Println("AL CONSUMIR AAAAAAAAAAAAAAAAA ->")

	incDto := &events.CreateEventDto{
		Type: events.Reposition,
		RepositionEvent: &events.RepositionEvent{
			ArticleId: message.Message.ArticleId,
			Quantity:  message.Message.Amount,
		},
	}

	if _, err := events.CreateEvent(incDto); err != nil {
		fmt.Println("ERROR AL CREAR EVENTO", err)
		return
	}

	sv, err := stockviews.GenerateStockView(message.Message.ArticleId)
	if err != nil {

		fmt.Println("ERROR AL REGENERAR STOCKVIEWS", message.Message.ArticleId)
	}

	if sv.Stock == message.Message.Amount { //Se repuso desde el 0 con el stock del ultimo msj
		emitStockNowAvailable(message.Message.ArticleId)
	}
}
