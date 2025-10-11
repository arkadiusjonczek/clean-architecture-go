package inmemory

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/entities"
)

var _ entities.BasketRepository = (*InMemoryBasketRepository)(nil)

type InMemoryBasketRepository struct {
	baskets map[string]*entities.Basket
}

func NewInMemoryBasketRepository() entities.BasketRepository {
	return &InMemoryBasketRepository{
		baskets: make(map[string]*entities.Basket),
	}
}

func (repository *InMemoryBasketRepository) Save(basket *entities.Basket) (string, error) {
	if basket == nil {
		return "", fmt.Errorf("basket is nil")
	}

	if basket.GetID() == "" {
		basket.SetID(uuid.NewString())
	}

	repository.baskets[basket.GetID()] = basket

	return basket.GetID(), nil
}

func (repository *InMemoryBasketRepository) Find(id string) (*entities.Basket, error) {
	basket, basketExists := repository.baskets[id]
	if !basketExists {
		return nil, fmt.Errorf("basket not found")
	}

	return basket, nil
}

func (repository *InMemoryBasketRepository) FindByUserId(userId string) (*entities.Basket, error) {
	for _, basket := range repository.baskets {
		if basket.GetUserID() == userId {
			return basket, nil
		}
	}

	return nil, &entities.BasketNotFoundError{}
}
