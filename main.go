package main

import (
	"emmanuel-guerreiro/stockgo/lib/db"
	"emmanuel-guerreiro/stockgo/rabbit"

	"emmanuel-guerreiro/stockgo/rest"
)

func main() {
	db.ConnectDatabase()
	defer db.DisconnectDatabase()

	rabbit.Init()
	rest.Init()
}
