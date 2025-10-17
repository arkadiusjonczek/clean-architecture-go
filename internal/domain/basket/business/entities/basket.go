package entities

import "fmt"

type Basket struct {
	Id     string
	UserID string
	Items  map[string]*BasketItem
}

type BasketItem struct {
	ProductID string
	Count     int
}

func (basket *Basket) GetID() string {
	return basket.Id
}

func (basket *Basket) SetID(id string) {
	basket.Id = id
}

func (basket *Basket) GetUserID() string {
	return basket.UserID
}

func (basket *Basket) GetItems() map[string]*BasketItem {
	return basket.Items
}

func (basket *Basket) HasItem(productID string) bool {
	_, basketHasItem := basket.Items[productID]

	return basketHasItem
}

func (basket *Basket) GetItem(productID string) (*BasketItem, error) {
	if !basket.HasItem(productID) {
		return nil, fmt.Errorf("basket does not have item with id: %s", productID)
	}

	return basket.Items[productID], nil
}

func (basket *Basket) AddItem(productID string, count int) *BasketItem {
	basketItem, basketHasItem := basket.Items[productID]
	if basketHasItem {
		basketItem.Count += count
	} else {
		basketItem = &BasketItem{
			ProductID: productID,
			Count:     count,
		}

		basket.Items[productID] = basketItem
	}

	return basketItem
}

func (basket *Basket) RemoveItem(productID string) error {
	if !basket.HasItem(productID) {
		return fmt.Errorf("basket does not have item with id: %s", productID)
	}

	delete(basket.Items, productID)

	return nil
}

func (basket *Basket) Clear() {
	if len(basket.Items) > 0 {
		basket.Items = map[string]*BasketItem{}
	}
}

func (basketItem *BasketItem) GetProductID() string {
	return basketItem.ProductID
}

func (basketItem *BasketItem) GetCount() int {
	return basketItem.Count
}

func (basketItem *BasketItem) SetCount(count int) {
	basketItem.Count = count
}
