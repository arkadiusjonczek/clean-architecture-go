package entities

type Basket struct {
	ID     string
	UserID string
	Items  map[string]*BasketItem
}

type BasketItem struct {
	ProductID string
	Count     int
}
