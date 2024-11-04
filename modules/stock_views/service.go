package stockviews

import (
	"emmanuel-guerreiro/stockgo/modules/events"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func FindOneById(id string) (*StockView, error) {
	res, err := findByID(id)

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

func GenerateStockView(id string) (*StockView, error) {
	//Will query orders_placed, generate the view and save it before returning it
	stock, err := events.FindArticleStockFromEvents(id)
	if err != nil {
		return nil, err
	}

	createDto := articleStockDtoToCreateStockViewDto(stock)
	_id, err := CreateOne(createDto)
	if err != nil {
		return nil, err
	}

	return findByID(_id)
}

func articleStockDtoToCreateStockViewDto(dto *events.ArticleStockDto) *CreateStockViewDto {
	return &CreateStockViewDto{
		ArticleId: dto.ArticleId,
		Stock:     dto.Stock,
	}
}
