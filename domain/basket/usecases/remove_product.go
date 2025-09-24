package usecases

import (
	"fmt"

	"github.com/arkadiusjonczek/clean-architecture-go/domain/basket"
	"github.com/arkadiusjonczek/clean-architecture-go/domain/basket/usecases/helper"
	"github.com/arkadiusjonczek/clean-architecture-go/domain/warehouse"
)

type RemoveProductUseCaseInput struct {
	UserID    string
	ProductID string
	Count     int
}

type RemoveProductUseCaseOutput struct {
	UserBasket *basket.Basket
	Actions    map[string]string
}

type RemoveProductUseCase interface {
	Execute(input *RemoveProductUseCaseInput) (*RemoveProductUseCaseOutput, error)
}

func NewRemoveProductUseCaseImpl(basketService helper.BasketCreatorService, basketRepository basket.BasketRepository, productRepository warehouse.ProductRepository) *RemoveProductUseCaseImpl {
	return &RemoveProductUseCaseImpl{
		basketService:     basketService,
		basketRepository:  basketRepository,
		productRepository: productRepository,
	}
}

var _ RemoveProductUseCase = (*RemoveProductUseCaseImpl)(nil)

type RemoveProductUseCaseImpl struct {
	basketService     helper.BasketCreatorService
	basketRepository  basket.BasketRepository
	productRepository warehouse.ProductRepository
}

func (useCase *RemoveProductUseCaseImpl) validate(input *RemoveProductUseCaseInput) error {
	if input == nil {
		return fmt.Errorf("input is nil")
	} else if input.UserID == "" {
		return fmt.Errorf("UserID is empty")
	} else if input.ProductID == "" {
		return fmt.Errorf("ProductID is empty")
	} else if input.Count <= 0 {
		return fmt.Errorf("Count is invalid (must be greater than 0)")
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

	basketItem, basketItemExists := userBasket.Items[input.ProductID]
	if basketItemExists {
		basketItem.Count -= input.Count

		if basketItem.Count <= 0 {
			delete(userBasket.Items, input.ProductID)
		}

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
