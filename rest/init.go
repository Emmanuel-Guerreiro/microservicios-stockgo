package rest

import (
	"emmanuel-guerreiro/stockgo/lib"
	"fmt"
)

func Init() {
	createServer()
	initRouter()
	getServer().Run(fmt.Sprintf(":%d", lib.GetEnv().Port))
}
