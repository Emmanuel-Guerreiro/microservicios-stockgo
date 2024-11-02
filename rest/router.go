package rest

import artconfig "emmanuel-guerreiro/stockgo/modules/article_config"

func initRouter() {
	if server == nil {
		panic("Server non existant")
	}

	//v1 routes
	v1 := server.Group("/v1")
	artconfig.InitController(v1)
}
