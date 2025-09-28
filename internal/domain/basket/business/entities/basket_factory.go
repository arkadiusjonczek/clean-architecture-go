package entities

import "fmt"

type BasketFactory interface {
	NewBasket(UserID string) (*Basket, error)
}

var _ BasketFactory = (*BasketFactoryImpl)(nil)

type BasketFactoryImpl struct {
}

func NewBasketFactory() BasketFactory {
	return &BasketFactoryImpl{}
}

func (factory *BasketFactoryImpl) NewBasket(UserID string) (*Basket, error) {
	if UserID == "" {
		return nil, fmt.Errorf("UserID cannot be empty")
	}

	return &Basket{
		UserID: UserID,
		Items:  make(map[string]*BasketItem),
	}, nil
}
