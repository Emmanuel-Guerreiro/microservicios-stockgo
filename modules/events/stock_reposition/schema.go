package stockreposition

type consumeStockRepositionDto struct {
	CorrelationId string                            `json:"correlation_id" example:"123123" `
	RoutingKey    string                            `json:"routing_key" example:"Remote RoutingKey to Reply"`
	Exchange      string                            `json:"exchange" example:"Remote Exchange to Reply"`
	Message       *consumeStockRepositionMessageDto `json:"message"`
}

type consumeStockRepositionMessageDto struct {
	ArticleId string `json:"articleId" example:"ArticleId"`
	Amount    int    `json:"amount"`
}

// Si el consume de stock reposition resulta en que un producto pase de no tener stock a tener
// Se envia el siguiente evento al queue de stock available
type placeStockAvailableDto struct {
	CorrelationId string                         `json:"correlation_id" example:"123123" `
	Message       *placeStockAvailableMessageDto `json:"message"`
}

type placeStockAvailableMessageDto struct {
	ArticleId string `json:"articleId" example:"ArticleId"`
	Amount    int    `json:"amount"`
}

// {
// 	"message":{
// 		"articleId": "123",
// 		"amount": 10
// 	}
// }
