package usecases

import (
	"fmt"
	"log"

	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/entities"
	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/usecases/helper"
)

type ShowBasketUseCaseInput struct {
	UserID string
}

type ShowBasketUseCaseOutput struct {
	UserBasket *entities.Basket
}

type ShowBasketUseCase interface {
	Execute(input *ShowBasketUseCaseInput) (*ShowBasketUseCaseOutput, error)
}

func NewShowBasketUseCaseImpl(basketService helper.BasketCreatorService) ShowBasketUseCase {
	return &ShowBasketUseCaseImpl{
		basketService: basketService,
	}
}

var _ ShowBasketUseCase = (*ShowBasketUseCaseImpl)(nil)

type ShowBasketUseCaseImpl struct {
	basketService helper.BasketCreatorService
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
		return nil, err
	}

	output := &ShowBasketUseCaseOutput{
		UserBasket: userBasket,
	}

	log.Printf("output: %v", output)

	return output, nil
}
