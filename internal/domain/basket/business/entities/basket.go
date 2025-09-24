package entities

import (
	warehouse "github.com/arkadiusjonczek/clean-architecture-go/internal/domain/warehouse/business/entities"
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
