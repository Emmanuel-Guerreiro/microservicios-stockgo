package ordersplaced

type ConsumeOrderPlacedDto struct {
	CorrelationId string                        `json:"correlation_id" example:"123123" `
	Message       *ConsumeOrderPlacedMessageDto `json:"message"`
}

type ConsumeOrderPlacedMessageDto struct {
	OrderId  string                          `json:"orderId"`
	CartId   string                          `json:"cartId"`
	Articles []*ConsumeOrderPlacedArticleDto `json:"articles"`
}

type ConsumeOrderPlacedArticleDto struct {
	ArticleId string `json:"articleId"`
	Quantity  int    `json:"quantity"`
}

// {
// 	"correlation_id": "123123",
// 	"message": {
// 		"orderId": "123123",
// 		"cartId": "123123",
// 		"articles": [
// 			{
// 				"articleId": "123",
// 				"quantity": 1
// 			}
// 		]
// 	}
// }
