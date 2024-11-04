package events

func FindArticleStockFromEvents(articleId string) (*ArticleStockDto, error) {

	articleStock, err := findArticleStockById(articleId)
	if err != nil {
		return nil, err
	}

	return articleStock, nil
}

func CreateEvent(event *CreateEventDto) (*Event, error) {
	id, err := create(event)
	if err != nil {
		return nil, err
	}

	return findEventById(id)
}

func findEventById(id string) (*Event, error) {
	return findById(id)
}
