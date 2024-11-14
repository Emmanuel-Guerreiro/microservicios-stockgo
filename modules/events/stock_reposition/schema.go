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

type placeStockAvailableMessageDto struct {
	ArticleId string `json:"articleId" example:"ArticleId"`
}

// {
// 	"message":{
// 		"articleId": "123",
// 		"amount": 10
// 	}
// }
