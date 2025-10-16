package entities

import "fmt"

type BasketFactory interface {
	NewBasket(UserID string) (*Basket, error)
	NewBasketWithID(basketID string, userID string) (*Basket, error)
}

var _ BasketFactory = (*BasketFactoryImpl)(nil)

type BasketFactoryImpl struct {
}

func NewBasketFactory() BasketFactory {
	return &BasketFactoryImpl{}
}

func (factory *BasketFactoryImpl) NewBasket(userID string) (*Basket, error) {
	if userID == "" {
		return nil, fmt.Errorf("userID cannot be empty")
	}

	return &Basket{
		userID: userID,
		items:  map[string]*BasketItem{},
	}, nil
}

func (factory *BasketFactoryImpl) NewBasketWithID(basketID string, userID string) (*Basket, error) {
	if basketID == "" {
		return nil, fmt.Errorf("basketID cannot be empty")
	} else if userID == "" {
		return nil, fmt.Errorf("userID cannot be empty")
	}

	return &Basket{
		id:     basketID,
		userID: userID,
		items:  map[string]*BasketItem{},
	}, nil
}
