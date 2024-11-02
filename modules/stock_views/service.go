package stockviews

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func findOneById(id string, ctx context.Context) (*StockView, error) {
	res, err := findByID(id)

	if err != nil {

		if err == mongo.ErrNoDocuments {
			return GenerateStockView(id)
		}

		return nil, err
	}

	return res, nil
}

func GenerateStockView(id string) (*StockView, error) {
	//Will query orders_placed, generate the view and save it before returning it
	return nil, nil
}
