package entities

import "fmt"

type Basket struct {
	id     string
	userID string
	items  map[string]*BasketItem
}
type BasketItem struct {
	productID string
	count     int
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
		basketItem.count += count
	} else {
		basketItem = &BasketItem{
			productID: productID,
			count:     count,
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

func (basketItem *BasketItem) GetProductID() string {
	return basketItem.productID
}

func (basketItem *BasketItem) GetCount() int {
	return basketItem.count
}

func (basketItem *BasketItem) SetCount(count int) {
	basketItem.count = count
}
