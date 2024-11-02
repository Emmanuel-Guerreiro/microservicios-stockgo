package main

import (
	"emmanuel-guerreiro/stockgo/lib/db"
	rabbit "emmanuel-guerreiro/stockgo/rabbit/consume"
	"emmanuel-guerreiro/stockgo/rest"
)

func main() {
	db.ConnectDatabase()
	defer db.DisconnectDatabase()

	rabbit.Init()
	rest.Init()
}
