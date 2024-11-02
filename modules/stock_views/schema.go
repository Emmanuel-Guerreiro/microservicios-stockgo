package stockviews

import "time"

type StockView struct {
	ArticleId string    `bson:"articleId" json:"articleId" validate:"required,min=1,max=100"`
	Stock     int       `json:"stock"`
	Created   time.Time `bson:"createdAt" json:"createdAt"`
	Updated   time.Time `bson:"updatedAt" json:"updatedAt"`
}

type CreateStockViewDto struct {
	ArticleId string    `bson:"articleId" json:"articleId" validate:"required,min=1,max=100"`
	Stock     int       `json:"stock"`
	Created   time.Time `bson:"createdAt" json:"createdAt"`
	Updated   time.Time `bson:"updatedAt" json:"updatedAt"`
}

type ReplaceStockViewDto struct {
	Stock   int       `json:"stock"`
	Updated time.Time `json:"updatedAt"`
}
