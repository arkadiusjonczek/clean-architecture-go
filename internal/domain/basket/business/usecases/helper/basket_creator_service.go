package helper

import (
	"errors"

	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/entities"
)

// BasketCreatorService used for special business logic to prevent duplicate code
type BasketCreatorService interface {
	FindOrCreate(userID string) (*entities.Basket, error)
}

var _ BasketCreatorService = (*BasketCreatorServiceImpl)(nil)

type BasketCreatorServiceImpl struct {
	basketFactory    entities.BasketFactory
	basketRepository entities.BasketRepository
}

func NewBasketCreatorServiceImpl(basketFactory entities.BasketFactory, basketRepository entities.BasketRepository) BasketCreatorService {
	return &BasketCreatorServiceImpl{
		basketFactory:    basketFactory,
		basketRepository: basketRepository,
	}
}

// FindOrCreate will create a new basket if a basket could not be found for the given userID
func (service *BasketCreatorServiceImpl) FindOrCreate(userID string) (*entities.Basket, error) {
	userBasket, basketRepositoryErr := service.basketRepository.FindByUserId(userID)
	if basketRepositoryErr != nil {
		// if the user has no basket yet, create it
		if errors.Is(basketRepositoryErr, &entities.BasketNotFoundError{}) {
			basket, newBasketErr := service.basketFactory.NewBasket(userID)
			if newBasketErr != nil {
				return nil, newBasketErr
			}
			userBasketID, saveBasketErr := service.basketRepository.Save(basket)
			if saveBasketErr != nil {
				return nil, saveBasketErr
			}
			userBasket, basketRepositoryErr = service.basketRepository.Find(userBasketID)
			if basketRepositoryErr != nil {
				return nil, basketRepositoryErr
			}
		} else {
			return nil, basketRepositoryErr
		}
	}

	return userBasket, nil
}
