package artconfig

import (
	"context"
	"emmanuel-guerreiro/stockgo/lib"
	"time"
)

func FindOneById(id string, ctx ...interface{}) (*ArticleConfig, error) {
	return findByID(id, ctx)
}

func FindOrCreateDefault(id string, ctx context.Context) (*ArticleConfig, error) {
	articleConfig, err := findByArticleId(id)
	if err != nil {
		return nil, err
	}
	if articleConfig == nil {
		articleConfig, err = createDefaultArticleConfig(id, ctx)
		if err != nil {
			return nil, err
		}
	}

	return articleConfig, nil
}

func findAll(page int, size int, ctx context.Context) ([]ArticleConfig, error) {
	pagination := lib.GetPagination(page, size)

	return findAllPaginated(pagination, ctx)
}

func createArticleConfig(articleConfig *CreateArticleConfigDto, ctx context.Context) (*ArticleConfig, error) {

	id, err := create(articleConfig, ctx)

	if err != nil {
		return nil, err
	}

	return FindOneById(id, ctx)
}

func replaceArticleConfig(id string, articleConfig *ReplaceArticleConfigDto, ctx context.Context) (*ArticleConfig, error) {
	res, err := updateOne(id, articleConfig, ctx)
	if err != nil {
		return nil, err
	}

	if res == 0 {
		return nil, ErrID
	}

	return FindOneById(id, ctx)
}

func deleteArticleConfig(id string, ctx context.Context) (*ArticleConfig, error) {
	art, err := findByID(id, ctx)
	if err != nil {
		return nil, err
	}

	if art == nil {
		return nil, ErrID
	}

	err = deleteByID(id, ctx)
	if err != nil {
		return nil, err
	}

	return art, nil
}

// TODO: Implement a general level project configuration
// and fetch from there
func createDefaultArticleConfig(id string, ctx context.Context) (*ArticleConfig, error) {
	articleConfig := &CreateArticleConfigDto{
		ArticleId:        id,
		AlertMinQuantity: 1,
	}
	return createArticleConfig(articleConfig, ctx)
}

func createDtoToArticleConfig(articleConfig *CreateArticleConfigDto) *ArticleConfig {
	return &ArticleConfig{
		ArticleId:        articleConfig.ArticleId,
		AlertMinQuantity: articleConfig.AlertMinQuantity,
		Created:          time.Now(),
		Updated:          time.Now(),
	}
}

func replaceDtoToArticleConfig(articleConfig *ReplaceArticleConfigDto) *ArticleConfig {
	return &ArticleConfig{
		AlertMinQuantity: articleConfig.AlertMinQuantity,
		Updated:          time.Now(),
	}
}
