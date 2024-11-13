package stockviews

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

var collection *mongo.Collection

func createIndex(collection *mongo.Collection) error {
	indexModel := mongo.IndexModel{
		//ArticleId es el ID externo -> Las lecturas y patches se realizan sobre este
		Keys: bson.D{{Key: "articleId", Value: 1}, {Key: "_id", Value: 1}}, // 1 for ascending order
	}
	_, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
	return err
}

func dbCollection() *mongo.Collection {

	if collection == nil {
		database := db.Get()
		collection = database.Collection("stock_views")
		createIndex(collection)
	}

	return collection
}

func findByArticleId(articleId string) (*StockView, error) {
	var articleConfig StockView

	if err := dbCollection().FindOne(context.TODO(), bson.M{"articleId": articleId}).Decode(&articleConfig); err != nil {
		return nil, err
	}

	return &articleConfig, nil
}

func create(stockViewDto *CreateStockViewDto) (string, error) {

	result, err := dbCollection().InsertOne(context.TODO(), createDtoToStockView(stockViewDto))
	if err != nil {
		return "", err
	}

	if oid, ok := result.InsertedID.(bson.ObjectID); ok {
		return oid.Hex(), nil
	}

	return "", ErrID
}

func updateOrCreate(stockViewDto *CreateStockViewDto) (*StockView, error) {
	var updated StockView

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After).SetUpsert(true)
	filter := bson.M{"articleId": stockViewDto.ArticleId}
	update := bson.M{"$set": createDtoToStockView(stockViewDto)}

	err := dbCollection().FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&updated)
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func createDtoToStockView(articleConfig *CreateStockViewDto) *StockView {
	return &StockView{
		ArticleId: articleConfig.ArticleId,
		Stock:     articleConfig.Stock,
		Created:   time.Now(),
		Updated:   time.Now(),
	}
}
