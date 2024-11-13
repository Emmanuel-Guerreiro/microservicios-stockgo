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
	Snapshot   EventType = "stock_snapshot"
)

type EventStatus string

const (
	Success        EventStatus = "success"
	NotEnoughStock EventStatus = "not_enough_stock"
)

type CreateEventDto struct {
	Type            EventType        `bson:"type" validate:"required"`
	DecrementEvent  *DecrementEvent  `bson:"decrementEvent"`
	RepositionEvent *RepositionEvent `bson:"repositionEvent"`
	EventStatus     EventStatus      `bson:"status"`
}

// Estuctura basica de del evento
type Event struct {
	ID              bson.ObjectID    `bson:"_id,omitempty"`
	Type            EventType        `bson:"type" validate:"required"`
	DecrementEvent  *DecrementEvent  `bson:"decrementEvent"`
	RepositionEvent *RepositionEvent `bson:"repositionEvent"`
	SnapshotEvent   *SnapshotEvent   `bson:"snapshotEvent"`
	Created         time.Time        `bson:"created"`
	EventStatus     EventStatus      `bson:"status"`
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

type SnapshotEvent struct {
	ArticleId string `bson:"articleId" json:"articleId" validate:"required,min=1,max=100"`
	Quantity  int    `bson:"quantity" json:"quantity" validate:"required"`
}

type ArticleStockFromDbDto struct {
	Id        string `bson:"_id" json:"_id"`
	ArticleId string `bson:"articleId" json:"articleId" ` //Used for compatibility with the aggregate query
	Stock     int    `bson:"stock" json:"stock" validate:"required,min=1,max=100"`
}

type ArticleStockDto struct {
	ArticleId string `bson:"articleId" json:"articleId" validate:"required,min=1,max=100"`
	Stock     int    `bson:"stock" json:"stock" validate:"required,min=1,max=100"`
}
