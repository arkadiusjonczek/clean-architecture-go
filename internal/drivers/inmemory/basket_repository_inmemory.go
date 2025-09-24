package inmemory

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/arkadiusjonczek/clean-architecture-go/domain/basket"
)

var _ basket.BasketRepository = (*InMemoryBasketRepository)(nil)

type InMemoryBasketRepository struct {
	baskets map[string]*basket.Basket
}

func NewInMemoryBasketRepository() basket.BasketRepository {
	return &InMemoryBasketRepository{
		baskets: make(map[string]*basket.Basket),
	}
}

func (repository *InMemoryBasketRepository) Save(basket *basket.Basket) (string, error) {
	if basket == nil {
		return "", fmt.Errorf("basket is nil")
	}

	if basket.ID == "" {
		basket.ID = uuid.NewString()
	}

	repository.baskets[basket.ID] = basket

	return basket.ID, nil
}

func (repository *InMemoryBasketRepository) Find(id string) (*basket.Basket, error) {
	basket, basketExists := repository.baskets[id]
	if !basketExists {
		return nil, fmt.Errorf("basket not found")
	}

	return basket, nil
}

func (repository *InMemoryBasketRepository) FindByUserId(userId string) (*basket.Basket, error) {
	for _, basket := range repository.baskets {
		if basket.UserID == userId {
			return basket, nil
		}
	}

	return nil, &basket.BasketNotFoundError{}
}
