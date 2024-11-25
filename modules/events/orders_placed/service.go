package ordersplaced

import (
	"emmanuel-guerreiro/stockgo/modules/events"
	stockviews "emmanuel-guerreiro/stockgo/modules/stock_views"
	"encoding/json"
	"fmt"
	"sync"
)

type stockStatus struct {
	hasStock bool
	article  *ConsumeOrderPlacedArticleDto
}

func ProcessOrderPlaced(data *ConsumeOrderPlacedDto) {
	var wg1 sync.WaitGroup
	ch := make(chan *stockStatus)
	for _, article := range data.Message.Articles {
		wg1.Add(1)
		go func(article *ConsumeOrderPlacedArticleDto) {
			defer wg1.Done()
			v, err := stockviews.FindOneById(article.ArticleId)
			if err != nil || v == nil { //Puede pasar que no exista el article?
				ch <- &stockStatus{hasStock: false, article: nil}
			}
			art := article
			q := article.Quantity
			s := v.Stock
			hs := s > q
			ch <- &stockStatus{hasStock: hs, article: art}
		}(article)
	}

	go func() {
		wg1.Wait()
		close(ch)
	}()

	for stockStatus := range ch {
		if !stockStatus.hasStock {
			if err := emitNotEnoughStock(stockStatus.article.ArticleId, stockStatus.article.Quantity, data.CorrelationId); err != nil {
				fmt.Println("ERROR AL EMITER NOT ENOUGH STOCK", err)
				return
			}
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
				EventStatus: events.Success,
			}

			if _, err := events.CreateEvent(decDto); err != nil {
				fmt.Println("ERROR AL CREAR EVENTO", err)
				return
			}

			if _, err := stockviews.GenerateStockViewNotify(article.ArticleId, data.CorrelationId); err != nil {
			}

		}(article)
	}

	wg.Wait()

	bodyParsed, err := json.Marshal(data)
	if err != nil {
		fmt.Println("ERROR AL MARSHAL STOCKVIEWS", data.CorrelationId)
		return
	}
	fmt.Println("Processed order placed succesfully", string(bodyParsed))
}
