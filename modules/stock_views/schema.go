package stockviews

import "time"

type StockView struct {
	ArticleId string    `bson:"articleId" json:"articleId" validate:"required,min=1,max=100"`
	Stock     int       `bson:"stock" json:"stock"`
	Created   time.Time `bson:"createdAt" json:"createdAt"`
	Updated   time.Time `bson:"updatedAt" json:"updatedAt"`
}

type CreateStockViewDto struct {
	ArticleId string `bson:"articleId" json:"articleId" validate:"required,min=1,max=100"`
	Stock     int    `bson:"stock" json:"stock"`
}

type ReplaceStockViewDto struct {
	Stock   int       `json:"stock"`
	Updated time.Time `json:"updatedAt"`
}

type StockConsultEvent struct {
	ArticleId string `json:"articleId"`
}

type StockViewResponseDto struct {
	ArticleId string `json:"articleId"`
	Stock     int    `json:"stock"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type stockConsultDto struct {
	ArticleId string `json:"articleId"`
}
