package usecases

import (
	"fmt"

	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/entities"
	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/usecases/helper"
	warehouse "github.com/arkadiusjonczek/clean-architecture-go/internal/domain/warehouse/business/entities"
)

type RemoveProductUseCaseInput struct {
	UserID    string
	ProductID string
}

type RemoveProductUseCaseOutput struct {
	UserBasket *entities.Basket
	Actions    map[string]string
}

type RemoveProductUseCase interface {
	Execute(input *RemoveProductUseCaseInput) (*RemoveProductUseCaseOutput, error)
}

func NewRemoveProductUseCaseImpl(basketService helper.BasketCreatorService, basketRepository entities.BasketRepository, productRepository warehouse.ProductRepository) RemoveProductUseCase {
	return &RemoveProductUseCaseImpl{
		basketService:     basketService,
		basketRepository:  basketRepository,
		productRepository: productRepository,
	}
}

var _ RemoveProductUseCase = (*RemoveProductUseCaseImpl)(nil)

type RemoveProductUseCaseImpl struct {
	basketService     helper.BasketCreatorService
	basketRepository  entities.BasketRepository
	productRepository warehouse.ProductRepository
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

	//basketItemCount := input.Count
	//// TODO: check for product stock in warehouse
	//if product.Stock < basketItemCount {
	//	return nil, fmt.Errorf("product stock is not enough")
	//}

	_, basketItemExists := userBasket.Items[input.ProductID]
	if basketItemExists {
		delete(userBasket.Items, input.ProductID)

		_, basketRepositorySaveErr := useCase.basketRepository.Save(userBasket)
		if basketRepositorySaveErr != nil {
			return nil, basketRepositorySaveErr
		}
	}

	output := &RemoveProductUseCaseOutput{
		UserBasket: userBasket,
		Actions:    map[string]string{},
	}

	return output, nil
}
