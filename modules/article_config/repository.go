package artconfig

import (
	"context"
	"emmanuel-guerreiro/stockgo/lib"
	"emmanuel-guerreiro/stockgo/lib/db"
	"emmanuel-guerreiro/stockgo/lib/log"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var ErrID = lib.NewValidationError().Add("id", "Invalid")

var collection *mongo.Collection

func dbCollection() *mongo.Collection {

	if collection == nil {
		database := db.Get()
		collection = database.Collection("article_config")
	}

	return collection
}

func findByID(id string, ctx ...interface{}) (*ArticleConfig, error) {
	var articleConfig ArticleConfig

	_id, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Get(ctx...).Error(err)
		return nil, ErrID
	}

	if err = dbCollection().FindOne(context.TODO(), bson.M{"_id": _id}).Decode(&articleConfig); err != nil {
		return nil, err
	}

	return &articleConfig, nil
}

func findByArticleId(id string) (*ArticleConfig, error) {
	var articleConfig ArticleConfig

	if err := dbCollection().FindOne(context.TODO(), bson.M{"articleId": id}).Decode(&articleConfig); err != nil {
		return nil, err
	}

	return &articleConfig, nil
}

func findAllPaginated(pagination *lib.Pagination, ctx context.Context) ([]ArticleConfig, error) {
	var articleConfig []ArticleConfig
	optionsFind := options.Find().SetSkip(pagination.Skip).SetLimit(pagination.Limit)

	cursor, err := dbCollection().Find(ctx, bson.M{}, optionsFind)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &articleConfig); err != nil {
		return nil, err
	}

	return articleConfig, nil
}

func create(articleConfig *CreateArticleConfigDto, ctx context.Context) (id string, err error) {

	result, err := dbCollection().InsertOne(ctx, createDtoToArticleConfig(articleConfig))
	if err != nil {
		return "", err
	}

	if oid, ok := result.InsertedID.(bson.ObjectID); ok {
		return oid.Hex(), nil
	}

	return "", ErrID

}

func updateOne(id string, articleConfig *ReplaceArticleConfigDto, ctx context.Context) (int64, error) {
	_id, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return 0, ErrID
	}

	updateResult, err := dbCollection().UpdateOne(ctx, bson.M{"_id": _id}, bson.M{"$set": replaceDtoToArticleConfig(articleConfig)})
	if err != nil {
		return 0, err
	}
	return updateResult.ModifiedCount, nil
}

func deleteByID(id string, ctx context.Context) error {
	_id, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return ErrID
	}

	_, err = dbCollection().DeleteOne(ctx, bson.M{"_id": _id})
	if err != nil {
		return err
	}
	return nil
}
