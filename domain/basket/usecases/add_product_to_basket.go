package usecases

import (
	"errors"
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
	basketFactory     basket.BasketFactory
	basketRepository  basket.BasketRepository
	productRepository warehouse.ProductRepository
}

func (useCase *AddProductToBasketUseCaseImpl) validate(input *AddProductToBasketUseCaseInput) error {
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

func (useCase *AddProductToBasketUseCaseImpl) Execute(input *AddProductToBasketUseCaseInput) (*AddProductToBasketUseCaseOutput, error) {
	// validate input first
	err := useCase.validate(input)
	if err != nil {
		return nil, fmt.Errorf("input validation error: %w", err)
	}

	// find the basket of the user
	userBasket, basketRepositoryErr := useCase.basketRepository.FindByUserId(input.UserID)
	if basketRepositoryErr != nil {
		// if the user has no basket yet, create it
		if errors.Is(basketRepositoryErr, &basket.BasketNotFoundError{}) {
			basket, newBasketErr := useCase.basketFactory.NewBasket(input.UserID)
			if newBasketErr != nil {
				return nil, newBasketErr
			}
			userBasketID, saveBasketErr := useCase.basketRepository.Save(basket)
			if saveBasketErr != nil {
				return nil, saveBasketErr
			}
			userBasket, basketRepositoryErr = useCase.basketRepository.Find(userBasketID)
			if basketRepositoryErr != nil {
				return nil, basketRepositoryErr
			}
		}

		return nil, basketRepositoryErr
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
	if !basketItemExists {
		basketItem = &basket.BasketItem{
			Product: product,
			Count:   0,
		}

		userBasket.Items[input.ProductID] = basketItem
	}

	basketItem.Count += input.Count

	_, basketRepositorySaveErr := useCase.basketRepository.Save(userBasket)
	if basketRepositorySaveErr != nil {
		return nil, basketRepositorySaveErr
	}

	output := &AddProductToBasketUseCaseOutput{}

	return output, nil
}
