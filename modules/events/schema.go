package events

import (
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type EventType string

const (
	Decrement  EventType = "stock_decrement"
	Reposition EventType = "stock_reposition"
)

type CreateEventDto struct {
	Type            EventType        `bson:"type" validate:"required"`
	DecrementEvent  *DecrementEvent  `bson:"decrementEvent"`
	RepositionEvent *RepositionEvent `bson:"repositionEvent"`
}

// Estuctura basica de del evento
type Event struct {
	ID              bson.ObjectID    `bson:"_id,omitempty"`
	Type            EventType        `bson:"type" validate:"required"`
	DecrementEvent  *DecrementEvent  `bson:"decrementEvent"`
	RepositionEvent *RepositionEvent `bson:"repositionEvent"`
	Created         time.Time        `bson:"created"`
}

// ValidateSchema valida la estructura para ser insertada en la db
func (e *Event) ValidateSchema() error {
	return validator.New().Struct(e)
}

type DecrementEvent struct {
	ArticleId string `bson:"articleId" json:"articleId" validate:"required,min=1,max=100"`
	Quantity  int    `bson:"quantity" json:"quantity" validate:"required,min=1,max=100"`
}

type RepositionEvent struct {
	ArticleId string `bson:"articleId" json:"articleId" validate:"required,min=1,max=100"`
	Quantity  int    `bson:"quantity" json:"quantity" validate:"required,min=1,max=100"`
}

type ArticleStockDto struct {
	ArticleId string `bson:"articleId" json:"articleId" validate:"required,min=1,max=100"`
	Stock     int    `bson:"stock" json:"stock" validate:"required,min=1,max=100"`
}
