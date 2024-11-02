package ordersplaced

import (
	stockviews "emmanuel-guerreiro/stockgo/modules/stock_views"
	"fmt"
	"sync"
)

func ProcessOrderPlaced(data *ConsumeOrderPlacedDto) {
	fmt.Println("WILL PROCESS ORDER")

	//Regenero la lista para todos los articulos incluidos en el pedido
	var wg sync.WaitGroup
	wg.Add(len(data.Message.Articles))
	for _, article := range data.Message.Articles {
		go func(a *ConsumeOrderPlacedArticleDto) {
			stockviews.GenerateStockView(a.ArticleId)
			wg.Done()
		}(article)
	}
	wg.Wait()
	return
}
