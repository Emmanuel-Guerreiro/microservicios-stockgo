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
