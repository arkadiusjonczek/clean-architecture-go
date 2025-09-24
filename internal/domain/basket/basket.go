package basket

import (
	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/warehouse"
)

type Basket struct {
	ID     string
	UserID string
	Items  map[string]*BasketItem
}

type BasketItem struct {
	Product *warehouse.Product
	Count   int
}
