package usecases

import (
	"fmt"
	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/usecases/dto"
	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/usecases/helper"
)

type ShowBasketUseCaseInput struct {
	UserID string
}

type ShowBasketUseCaseOutput struct {
	UserBasket *dto.BasketDTO
}

type ShowBasketUseCase interface {
	Execute(input *ShowBasketUseCaseInput) (*ShowBasketUseCaseOutput, error)
}

func NewShowBasketUseCaseImpl(basketService helper.BasketCreatorService, basketOutputService helper.BasketOutputService) ShowBasketUseCase {
	return &ShowBasketUseCaseImpl{
		basketService:       basketService,
		basketOutputService: basketOutputService,
	}
}

var _ ShowBasketUseCase = (*ShowBasketUseCaseImpl)(nil)

type ShowBasketUseCaseImpl struct {
	basketService       helper.BasketCreatorService
	basketOutputService helper.BasketOutputService
}

func (useCase *ShowBasketUseCaseImpl) validate(input *ShowBasketUseCaseInput) error {
	if input == nil {
		return fmt.Errorf("input is nil")
	} else if input.UserID == "" {
		return fmt.Errorf("UserID is empty")
	}

	return nil
}

func (useCase *ShowBasketUseCaseImpl) Execute(input *ShowBasketUseCaseInput) (*ShowBasketUseCaseOutput, error) {
	err := useCase.validate(input)
	if err != nil {
		return nil, fmt.Errorf("input validation error: %w", err)
	}

	userBasket, userBasketErr := useCase.basketService.FindOrCreate(input.UserID)
	if userBasketErr != nil {
		return nil, userBasketErr
	}

	userBasketDTO, basketOutputServiceErr := useCase.basketOutputService.CreateBasketDTO(userBasket)
	if basketOutputServiceErr != nil {
		return nil, basketOutputServiceErr
	}

	output := &ShowBasketUseCaseOutput{
		UserBasket: userBasketDTO,
	}

	return output, nil
}
