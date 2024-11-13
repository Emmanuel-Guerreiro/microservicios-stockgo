package artconfig

import "time"

type ArticleConfig struct {
	ArticleId        string    `bson:"articleId" json:"articleId" validate:"required,min=1,max=100"`
	AlertMinQuantity int       `bson:"alertMinQuantity" json:"alertMinQuantity" validate:"required,min=1,max=100"`
	Created          time.Time `bson:"createdAt" json:"createdAt"`
	Updated          time.Time `bson:"updatedAt" json:"updatedAt"`
}

type CreateArticleConfigDto struct {
	ArticleId        string `bson:"articleId" json:"articleId" validate:"required,min=1,max=100"`
	AlertMinQuantity int    `bson:"alertMinQuantity" json:"alertMinQuantity" validate:"required,min=1,max=100"`
}

type ReplaceArticleConfigDto struct {
	AlertMinQuantity int `bson:"alertMinQuantity" json:"alertMinQuantity" validate:"required,min=1,max=100"`
}

type ArticleFindResponsePaginated struct {
	Status int             `json:"status"`
	Data   []ArticleConfig `json:"data"`
	Page   int             `json:"page"`
	Length int             `json:"length"`
}
