package usecases

import (
	"fmt"

	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/entities"
	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/usecases/helper"
)

type ClearBasketUseCaseInput struct {
	UserID    string
	ProductID string
	Count     int
}

type ClearBasketUseCaseOutput struct {
	UserBasket *entities.Basket
}

type ClearBasketUseCase interface {
	Execute(input *ClearBasketUseCaseInput) (*ClearBasketUseCaseOutput, error)
}

func NewClearBasketUseCaseImpl(basketService helper.BasketCreatorService, basketRepository entities.BasketRepository) ClearBasketUseCase {
	return &ClearBasketUseCaseImpl{
		basketService:    basketService,
		basketRepository: basketRepository,
	}
}

var _ ClearBasketUseCase = (*ClearBasketUseCaseImpl)(nil)

type ClearBasketUseCaseImpl struct {
	basketService    helper.BasketCreatorService
	basketRepository entities.BasketRepository
}

func (useCase *ClearBasketUseCaseImpl) validate(input *ClearBasketUseCaseInput) error {
	if input == nil {
		return fmt.Errorf("input is nil")
	} else if input.UserID == "" {
		return fmt.Errorf("UserID is empty")
	}

	return nil
}

func (useCase *ClearBasketUseCaseImpl) Execute(input *ClearBasketUseCaseInput) (*ClearBasketUseCaseOutput, error) {
	err := useCase.validate(input)
	if err != nil {
		return nil, fmt.Errorf("input validation error: %w", err)
	}

	userBasket, userBasketErr := useCase.basketService.FindOrCreate(input.UserID)
	if userBasketErr != nil {
		return nil, err
	}

	// clear the basket by replacing the map
	userBasket.Items = make(map[string]*entities.BasketItem)

	_, basketRepositorySaveErr := useCase.basketRepository.Save(userBasket)
	if basketRepositorySaveErr != nil {
		return nil, basketRepositorySaveErr
	}

	output := &ClearBasketUseCaseOutput{
		UserBasket: userBasket,
	}

	return output, nil
}
