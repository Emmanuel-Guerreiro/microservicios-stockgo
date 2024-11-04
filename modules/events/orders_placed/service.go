package ordersplaced

import (
	"emmanuel-guerreiro/stockgo/modules/events"
	stockviews "emmanuel-guerreiro/stockgo/modules/stock_views"
	"fmt"
	"sync"
)

func ProcessOrderPlaced(data *ConsumeOrderPlacedDto) {
	fmt.Println("WILL PROCESS ORDER")
	var wg1 sync.WaitGroup

	/*
		TODO: Should be an Article channel to facilitate which articles are not available?
		TODO: If the stock is below a configured amount should I use a distributed mutex in order to avoid illegal stock values?
	*/
	//If any of the articles cant provide the required stock
	//Will not process the order
	//And notify through the event bus
	ch := make(chan bool)
	for _, article := range data.Message.Articles {
		wg1.Add(1)
		go func(article *ConsumeOrderPlacedArticleDto) {
			defer wg1.Done()
			v, err := stockviews.FindOneById(article.ArticleId)
			if err != nil || v == nil { //Puede pasar que no exista el article?
				ch <- false
			}

			ch <- v.Stock > article.Quantity
		}(article)
	}

	go func() {
		wg1.Wait()
		close(ch)
	}()

	for hasStock := range ch {
		if !hasStock {
			//TODO: Place an event in the event bus notifying that the order cant be processed
			fmt.Println("Stock insufficient for one or more articles.")
			return
		}
	}
	//For each article, will decrement the stock by appending the new decrement stock events
	//And recalculate the stock view
	var wg sync.WaitGroup
	wg.Add(len(data.Message.Articles))
	for _, article := range data.Message.Articles {
		go func(article *ConsumeOrderPlacedArticleDto) {
			defer wg.Done()

			decDto := &events.CreateEventDto{
				Type: events.Decrement,
				DecrementEvent: &events.DecrementEvent{
					ArticleId: article.ArticleId,
					Quantity:  article.Quantity,
				},
			}

			if _, err := events.CreateEvent(decDto); err != nil {
				//TODO: Place an event in the event bus notifying that the order cant be processed
				fmt.Println("ERROR AL CREAR EVENTO", err)
				return
			}

			if _, err := stockviews.GenerateStockView(article.ArticleId); err != nil {
				//TODO:Should notify to somewhere?
				fmt.Println("ERROR AL REGENERAR STOCKVIEWS", article.ArticleId)
			}
		}(article)
	}
	wg.Wait()

	return
}
