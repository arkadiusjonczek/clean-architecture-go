package usecases

import (
	"fmt"

	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket"
	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/usecases/helper"
	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/warehouse"
)

type UpdateProductCountUseCaseInput struct {
	UserID    string
	ProductID string
	Count     int
}

type UpdateProductCountUseCaseOutput struct {
	UserBasket *basket.Basket
	Actions    map[string]string
}

type UpdateProductCountUseCase interface {
	Execute(input *UpdateProductCountUseCaseInput) (*UpdateProductCountUseCaseOutput, error)
}

func NewUpdateProductCountImpl(basketService helper.BasketCreatorService, basketRepository basket.BasketRepository, productRepository warehouse.ProductRepository) UpdateProductCountUseCase {
	return &UpdateProductCountUseCaseImpl{
		basketService:     basketService,
		basketRepository:  basketRepository,
		productRepository: productRepository,
	}
}

var _ UpdateProductCountUseCase = (*UpdateProductCountUseCaseImpl)(nil)

type UpdateProductCountUseCaseImpl struct {
	basketService     helper.BasketCreatorService
	basketRepository  basket.BasketRepository
	productRepository warehouse.ProductRepository
}

func (useCase *UpdateProductCountUseCaseImpl) validate(input *UpdateProductCountUseCaseInput) error {
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

func (useCase *UpdateProductCountUseCaseImpl) Execute(input *UpdateProductCountUseCaseInput) (*UpdateProductCountUseCaseOutput, error) {
	// validate input first
	err := useCase.validate(input)
	if err != nil {
		return nil, fmt.Errorf("input validation error: %w", err)
	}

	userBasket, userBasketErr := useCase.basketService.FindOrCreate(input.UserID)
	if userBasketErr != nil {
		return nil, err
	}

	product, productRepositoryErr := useCase.productRepository.Find(input.ProductID)
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
		basketItem.Count = input.Count
	} else {
		userBasket.Items[input.ProductID] = &basket.BasketItem{
			Product: product,
			Count:   input.Count,
		}
	}

	_, basketRepositorySaveErr := useCase.basketRepository.Save(userBasket)
	if basketRepositorySaveErr != nil {
		return nil, basketRepositorySaveErr
	}

	output := &UpdateProductCountUseCaseOutput{
		UserBasket: userBasket,
		Actions:    map[string]string{},
	}

	return output, nil
}
