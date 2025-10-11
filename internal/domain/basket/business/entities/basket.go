package entities

import "fmt"

type Basket struct {
	id     string
	userID string
	items  map[string]*BasketItem
}
type BasketItem struct {
	ProductID string
	Count     int
}

func (basket *Basket) GetID() string {
	return basket.id
}

func (basket *Basket) SetID(id string) {
	basket.id = id
}

func (basket *Basket) GetUserID() string {
	return basket.userID
}

func (basket *Basket) GetItems() map[string]*BasketItem {
	return basket.items
}

func (basket *Basket) HasItem(productID string) bool {
	_, basketHasItem := basket.items[productID]

	return basketHasItem
}

func (basket *Basket) GetItem(productID string) (*BasketItem, error) {
	if !basket.HasItem(productID) {
		return nil, fmt.Errorf("basket does not have item with id: %s", productID)
	}

	return basket.items[productID], nil
}

func (basket *Basket) AddItem(productID string, count int) *BasketItem {
	basketItem, basketHasItem := basket.items[productID]
	if basketHasItem {
		basketItem.Count += count
	} else {
		basketItem = &BasketItem{
			ProductID: productID,
			Count:     count,
		}

		basket.items[productID] = basketItem
	}

	return basketItem
}

func (basket *Basket) RemoveItem(productID string) error {
	if !basket.HasItem(productID) {
		return fmt.Errorf("basket does not have item with id: %s", productID)
	}

	delete(basket.items, productID)

	return nil
}

func (basket *Basket) Clear() {
	if len(basket.items) > 0 {
		basket.items = map[string]*BasketItem{}
	}
}
