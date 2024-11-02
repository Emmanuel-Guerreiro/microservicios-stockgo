package stockviews

import (
	"context"
	"emmanuel-guerreiro/stockgo/lib"
	"emmanuel-guerreiro/stockgo/lib/db"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var ErrID = lib.NewValidationError().Add("id", "Invalid")

var collection *mongo.Collection

func dbCollection() *mongo.Collection {

	if collection == nil {
		database := db.Get()
		collection = database.Collection("stock_views")
	}

	return collection
}

func findByID(id string) (*StockView, error) {
	var articleConfig StockView

	_id, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrID
	}

	if err = dbCollection().FindOne(context.TODO(), bson.M{"_id": _id}).Decode(&articleConfig); err != nil {
		return nil, err
	}

	return &articleConfig, nil
}
