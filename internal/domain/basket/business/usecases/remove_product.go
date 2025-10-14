package usecases

import (
	"fmt"

	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/entities"
	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/usecases/dto"
	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/usecases/helper"
	warehouse "github.com/arkadiusjonczek/clean-architecture-go/internal/domain/warehouse/business/entities"
)

type RemoveProductUseCaseInput struct {
	UserID    string
	ProductID string
}

type RemoveProductUseCaseOutput struct {
	UserBasket *dto.BasketDTO
	Actions    map[string]string
}

type RemoveProductUseCase interface {
	Execute(input *RemoveProductUseCaseInput) (*RemoveProductUseCaseOutput, error)
}

func NewRemoveProductUseCaseImpl(basketService helper.BasketCreatorService, basketOutputService helper.BasketOutputService, basketRepository entities.BasketRepository, productRepository warehouse.ProductRepository) RemoveProductUseCase {
	return &RemoveProductUseCaseImpl{
		basketService:       basketService,
		basketOutputService: basketOutputService,
		basketRepository:    basketRepository,
		productRepository:   productRepository,
	}
}

var _ RemoveProductUseCase = (*RemoveProductUseCaseImpl)(nil)

type RemoveProductUseCaseImpl struct {
	basketService       helper.BasketCreatorService
	basketOutputService helper.BasketOutputService
	basketRepository    entities.BasketRepository
	productRepository   warehouse.ProductRepository
}

func (useCase *RemoveProductUseCaseImpl) validate(input *RemoveProductUseCaseInput) error {
	if input == nil {
		return fmt.Errorf("input is nil")
	} else if input.UserID == "" {
		return fmt.Errorf("UserID is empty")
	} else if input.ProductID == "" {
		return fmt.Errorf("ProductID is empty")
	}

	return nil
}

func (useCase *RemoveProductUseCaseImpl) Execute(input *RemoveProductUseCaseInput) (*RemoveProductUseCaseOutput, error) {
	// validate input first
	err := useCase.validate(input)
	if err != nil {
		return nil, fmt.Errorf("input validation error: %w", err)
	}

	userBasket, userBasketErr := useCase.basketService.FindOrCreate(input.UserID)
	if userBasketErr != nil {
		return nil, err
	}

	_, productRepositoryErr := useCase.productRepository.Find(input.ProductID)
	if productRepositoryErr != nil {
		return nil, productRepositoryErr
	}

	_ = userBasket.RemoveItem(input.ProductID)

	_, basketRepositorySaveErr := useCase.basketRepository.Save(userBasket)
	if basketRepositorySaveErr != nil {
		return nil, basketRepositorySaveErr
	}

	userBasketDTO, basketOutputServiceErr := useCase.basketOutputService.CreateBasketDTO(userBasket)
	if basketOutputServiceErr != nil {
		return nil, basketOutputServiceErr
	}

	output := &RemoveProductUseCaseOutput{
		UserBasket: userBasketDTO,
		Actions:    map[string]string{},
	}

	return output, nil
}
