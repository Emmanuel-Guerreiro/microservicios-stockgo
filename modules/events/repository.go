package events

import (
	"context"
	"emmanuel-guerreiro/stockgo/lib"
	"emmanuel-guerreiro/stockgo/lib/db"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var ErrID = lib.NewValidationError().Add("id", "Invalid")
var ErrType = lib.NewValidationError().Add("type", "Invalid")

var collection *mongo.Collection

func createPartialIndex(collection *mongo.Collection) error {
	// Define the index model with the field you want to index.
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "type", Value: 1}, {Key: "snapshotEvent.articleId", Value: 1}}, // 1 for ascending order
		Options: options.Index().SetPartialFilterExpression(bson.D{
			{Key: "type", Value: Snapshot},
		}),
	}

	// Create the index on the collection
	_, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
	return err
}

func dbCollection() *mongo.Collection {

	if collection == nil {
		database := db.Get()
		collection = database.Collection("events")
		createPartialIndex(collection) //Genera el index para filtrar mas rapido los snapshots
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
	matchStage := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "$or", Value: bson.A{
				bson.D{{Key: "type", Value: "stock_reposition"}, {Key: "repositionEvent.articleId", Value: articleID}},
				bson.D{{Key: "type", Value: "stock_decrement"}, {Key: "decrementEvent.articleId", Value: articleID}},
			}},
		}},
	}
	cursor, err := dbCollection().Aggregate(context.TODO(), mongo.Pipeline{matchStage})
	if err != nil {
		return nil, err
	}

	var results []Event
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, nil
	}

	var totalStock int = 0
	for _, result := range results {
		if result.Type == Reposition {
			totalStock += result.RepositionEvent.Quantity
		} else if result.Type == Decrement {
			totalStock -= result.DecrementEvent.Quantity
		}
	}

	return &ArticleStockDto{
		ArticleId: articleID,
		Stock:     totalStock,
	}, nil
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
