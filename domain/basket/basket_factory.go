package basket

import "fmt"

func NewBasket(ID string, UserID string) (*Basket, error) {
	if ID == "" {
		return nil, fmt.Errorf("ID cannot be empty")
	} else if UserID == "" {
		return nil, fmt.Errorf("UserID cannot be empty")
	}

	return &Basket{
		ID:     ID,
		UserID: UserID,
		Items:  make(map[string]*BasketItem),
	}, nil
}
