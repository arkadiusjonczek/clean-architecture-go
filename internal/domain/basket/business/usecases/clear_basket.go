package usecases

import (
	"fmt"

	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/entities"
	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/usecases/dto"
	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/usecases/helper"
)

type ClearBasketUseCaseInput struct {
	UserID    string
	ProductID string
	Count     int
}

type ClearBasketUseCaseOutput struct {
	UserBasket *dto.BasketDTO
}

type ClearBasketUseCase interface {
	Execute(input *ClearBasketUseCaseInput) (*ClearBasketUseCaseOutput, error)
}

func NewClearBasketUseCaseImpl(basketService helper.BasketCreatorService, basketOutputService helper.BasketOutputService, basketRepository entities.BasketRepository) ClearBasketUseCase {
	return &ClearBasketUseCaseImpl{
		basketService:       basketService,
		basketOutputService: basketOutputService,
		basketRepository:    basketRepository,
	}
}

var _ ClearBasketUseCase = (*ClearBasketUseCaseImpl)(nil)

type ClearBasketUseCaseImpl struct {
	basketService       helper.BasketCreatorService
	basketOutputService helper.BasketOutputService
	basketRepository    entities.BasketRepository
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

	userBasket.Clear()

	_, basketRepositorySaveErr := useCase.basketRepository.Save(userBasket)
	if basketRepositorySaveErr != nil {
		return nil, basketRepositorySaveErr
	}

	userBasketDTO, basketOutputServiceErr := useCase.basketOutputService.CreateBasketDTO(userBasket)
	if basketOutputServiceErr != nil {
		return nil, basketOutputServiceErr
	}

	output := &ClearBasketUseCaseOutput{
		UserBasket: userBasketDTO,
	}

	return output, nil
}
