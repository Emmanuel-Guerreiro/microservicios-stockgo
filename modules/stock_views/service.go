package stockviews

import (
	"context"
	artconfig "emmanuel-guerreiro/stockgo/modules/article_config"
	"emmanuel-guerreiro/stockgo/modules/events"
	requirereposition "emmanuel-guerreiro/stockgo/modules/require_reposition"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func FindOneById(id string) (*StockView, error) {
	res, err := findByArticleId(id)

	if err != nil {

		if err != mongo.ErrNoDocuments {
			return nil, err
		}

		return GenerateStockView(id)
	}

	return res, nil
}

func CreateOne(dto *CreateStockViewDto) (string, error) {
	return create(dto)
}

func updateOrCreateOne(dto *CreateStockViewDto) (*StockView, error) {
	return updateOrCreate(dto)
}

func GenerateStockView(id string) (*StockView, error) {
	//Will query orders_placed, generate the view and save it before returning it
	stock, err := events.FindArticleStockFromEvents(id)
	if err != nil {
		return nil, err
	}

	if stock == nil {
		stock = &events.ArticleStockDto{
			ArticleId: id,
			Stock:     0,
		}
	}

	createDto := articleStockDtoToCreateStockViewDto(stock)
	sv, err := updateOrCreateOne(createDto)
	if err != nil {
		return nil, err
	}

	return sv, nil
}

func GenerateStockViewNotify(id string) (*StockView, error) {
	sv, err := GenerateStockView(id)
	if err != nil {
		return nil, err
	}

	config, err := artconfig.FindOrCreateDefault(id, context.TODO())
	if err != nil {
		return sv, err
	}

	if config.AlertMinQuantity <= sv.Stock {
		//Notifico que hace falta comprar
		requirereposition.EmitNotEnoughStock(id)
	}

	return sv, nil
}

func articleStockDtoToCreateStockViewDto(dto *events.ArticleStockDto) *CreateStockViewDto {
	return &CreateStockViewDto{
		ArticleId: dto.ArticleId,
		Stock:     dto.Stock,
	}
}
