package artconfig

import (
	"context"
	"time"
)

func findOneById(id string, ctx ...interface{}) (*ArticleConfig, error) {
	return findByID(id, ctx)
}

func createArticleConfig(articleConfig *CreateArticleConfigDto, ctx context.Context) (*ArticleConfig, error) {

	id, err := create(articleConfig, ctx)

	if err != nil {
		return nil, err
	}

	return findOneById(id, ctx)
}

func replaceArticleConfig(id string, articleConfig *ReplaceArticleConfigDto, ctx context.Context) (*ArticleConfig, error) {
	res, err := updateOne(id, articleConfig, ctx)
	if err != nil {
		return nil, err
	}

	if res == 0 {
		return nil, ErrID
	}

	return findOneById(id, ctx)
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
