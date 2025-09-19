package usecases

import (
	"fmt"

	"github.com/arkadiusjonczek/clean-architecture-go/domain/basket"
	"github.com/arkadiusjonczek/clean-architecture-go/domain/warehouse"
)

type AddProductToBasketUseCaseInput struct {
	UserID    string
	ProductID string
	Count     int
}

type AddProductToBasketUseCaseOutput struct {
	actions map[string]string
}

type AddProductToBasketUseCase interface {
	Execute(input *AddProductToBasketUseCaseInput) (*AddProductToBasketUseCaseOutput, error)
}

func NewAddProductToBasketUseCaseImpl(basketRepository basket.BasketRepository, productRepository warehouse.ProductRepository) *AddProductToBasketUseCaseImpl {
	return &AddProductToBasketUseCaseImpl{
		basketRepository:  basketRepository,
		productRepository: productRepository,
	}
}

var _ AddProductToBasketUseCase = (*AddProductToBasketUseCaseImpl)(nil)

type AddProductToBasketUseCaseImpl struct {
	basketRepository  basket.BasketRepository
	productRepository warehouse.ProductRepository
}

func (useCase *AddProductToBasketUseCaseImpl) Execute(input *AddProductToBasketUseCaseInput) (*AddProductToBasketUseCaseOutput, error) {
	// validate input first
	if input == nil {
		return nil, fmt.Errorf("input is nil")
	} else if input.UserID == "" {
		return nil, fmt.Errorf("UserID is empty")
	} else if input.ProductID == "" {
		return nil, fmt.Errorf("ProductID is empty")
	} else if input.Count <= 0 {
		return nil, fmt.Errorf("Count is invalid")
	}

	// TODO: how to handle guest and logged in customers?
	userBasket, basketRepositoryErr := useCase.basketRepository.FindByUserId(input.UserID)
	if basketRepositoryErr != nil {
		return nil, basketRepositoryErr
	}

	if basketItem, basketItemExists := userBasket.Items[input.ProductID]; basketItemExists {
		product, productRepositoryErr := useCase.productRepository.Find(input.ProductID)
		if productRepositoryErr != nil {
			return nil, productRepositoryErr
		}

		if product.Stock < basketItem.Count+input.Count {

		}

		basketItem.Count += input.Count
	} else {
		product, productRepositoryErr := useCase.productRepository.Find(input.ProductID)
		if productRepositoryErr != nil {
			return nil, productRepositoryErr
		}

		userBasket.Items[input.ProductID] = &basket.BasketItem{
			Product: *product,
			Count:   input.Count,
		}
	}

	_, basketRepositorySaveErr := useCase.basketRepository.Save(userBasket)
	if basketRepositorySaveErr != nil {
		return nil, basketRepositorySaveErr
	}

	output := &AddProductToBasketUseCaseOutput{}

	return output, nil
}
