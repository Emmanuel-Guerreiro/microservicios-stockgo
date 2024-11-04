package events

import (
	"context"
	"emmanuel-guerreiro/stockgo/lib"
	"emmanuel-guerreiro/stockgo/lib/db"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var ErrID = lib.NewValidationError().Add("id", "Invalid")
var ErrType = lib.NewValidationError().Add("type", "Invalid")

var collection *mongo.Collection

func dbCollection() *mongo.Collection {

	if collection == nil {
		database := db.Get()
		collection = database.Collection("article_config")
	}

	return collection
}

func findAllByType(eventType EventType) (*[]Event, error) {

	var event []Event

	cursor, err := dbCollection().Find(context.TODO(), bson.M{"type": eventType})
	if err != nil {
		return nil, err
	}

	if err := cursor.All(context.TODO(), &event); err != nil {
		return nil, err
	}

	return &event, nil
}

func findArticleStockById(articleID string) (*ArticleStockDto, error) {
	// create group stage
	_id, err := bson.ObjectIDFromHex(articleID)
	// Create match stage to filter events by article ID and relevant types
	// Create match stage to filter events by article ID and relevant types
	matchStage := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "article.id", Value: _id},
			{Key: "type", Value: bson.D{
				{Key: "$in", Value: bson.A{"stock_reposition", "stock_decrement"}},
			}},
		}},
	}

	// Create group stage to calculate the total amount
	groupStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$article.id"},
			{Key: "totalAmount", Value: bson.D{
				{Key: "$sum", Value: bson.D{
					{Key: "$cond", Value: bson.D{
						{Key: "if", Value: bson.D{{Key: "$eq", Value: bson.A{"$type", "stock_reposition"}}}},
						{Key: "then", Value: "$article.amount"},
						{Key: "else", Value: bson.D{{Key: "$multiply", Value: bson.A{"$article.amount", -1}}}},
					}},
				}},
			}},
		}},
	}

	// pass the pipeline to the Aggregate() method
	cursor, err := dbCollection().Aggregate(context.TODO(), mongo.Pipeline{matchStage, groupStage})
	if err != nil {
		return nil, err
	}

	var results []ArticleStockDto
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}
	return &results[0], nil
}

func create(event *CreateEventDto) (string, error) {
	e, err := createEventDtoToEvent(event)
	if err != nil {
		return "", err
	}

	result, err := dbCollection().InsertOne(context.TODO(), e)
	if err != nil {
		return "", err
	}

	if oid, ok := result.InsertedID.(bson.ObjectID); ok {
		return oid.Hex(), nil
	}

	return "", ErrID
}

func createEventDtoToEvent(event *CreateEventDto) (*Event, error) {
	switch event.Type {
	case Decrement:
		return &Event{
			Type:            Decrement,
			DecrementEvent:  event.DecrementEvent,
			RepositionEvent: nil,
			Created:         time.Now(),
		}, nil
	case Reposition:
		return &Event{
			Type:            Reposition,
			DecrementEvent:  nil,
			RepositionEvent: event.RepositionEvent,
			Created:         time.Now(),
		}, nil
	default:
		return nil, ErrType
	}
}

func findById(id string) (*Event, error) {
	var event Event
	_id, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrID
	}

	err = dbCollection().FindOne(context.TODO(), bson.M{"_id": _id}).Decode(&event)
	if err != nil {
		return nil, err
	}

	return &event, nil
}
