package entities

import "fmt"

type BasketFactory interface {
	NewBasket(UserID string) (*Basket, error)
	NewBasketWithID(userID string, id string) (*Basket, error)
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

func (factory *BasketFactoryImpl) NewBasketWithID(userID string, id string) (*Basket, error) {
	if userID == "" {
		return nil, fmt.Errorf("userID cannot be empty")
	} else if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	return &Basket{
		id:     id,
		userID: userID,
		items:  map[string]*BasketItem{},
	}, nil
}
