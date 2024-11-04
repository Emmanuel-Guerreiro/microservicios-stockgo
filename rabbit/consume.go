package rabbit

import (
	ordersplaced "emmanuel-guerreiro/stockgo/modules/events/orders_placed"
	stockreposition "emmanuel-guerreiro/stockgo/modules/events/stock_reposition"
	"emmanuel-guerreiro/stockgo/security"
)

func Init() {
	go stockreposition.ListenerReposition()
	go ordersplaced.ListenerOrderPlaced()
	go security.ListenerLogout()
}
