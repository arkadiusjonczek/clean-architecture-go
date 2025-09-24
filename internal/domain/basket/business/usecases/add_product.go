package usecases

import (
	"fmt"

	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/entities"
	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/usecases/dto"
	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/usecases/helper"
	warehouse "github.com/arkadiusjonczek/clean-architecture-go/internal/domain/warehouse/business/entities"
)

type AddProductUseCaseInput struct {
	UserID    string
	ProductID string
	Count     int
}

type AddProductUseCaseOutput struct {
	UserBasketDTO *dto.BasketDTO
	Actions       map[string]string
}

type AddProductUseCase interface {
	Execute(input *AddProductUseCaseInput) (*AddProductUseCaseOutput, error)
}

func NewAddProductUseCaseImpl(basketService helper.BasketCreatorService, basketOutputService helper.BasketOutputService, basketRepository entities.BasketRepository, productRepository warehouse.ProductRepository) AddProductUseCase {
	return &AddProductUseCaseImpl{
		basketService:       basketService,
		basketOutputService: basketOutputService,
		basketRepository:    basketRepository,
		productRepository:   productRepository,
	}
}

var _ AddProductUseCase = (*AddProductUseCaseImpl)(nil)

type AddProductUseCaseImpl struct {
	basketService       helper.BasketCreatorService
	basketOutputService helper.BasketOutputService
	basketRepository    entities.BasketRepository
	productRepository   warehouse.ProductRepository
}

func (useCase *AddProductUseCaseImpl) validate(input *AddProductUseCaseInput) error {
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

func (useCase *AddProductUseCaseImpl) Execute(input *AddProductUseCaseInput) (*AddProductUseCaseOutput, error) {
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
	if !basketItemExists {
		basketItem = &entities.BasketItem{
			Product:   product,
			ProductID: product.ID,
			Count:     0,
		}

		userBasket.Items[input.ProductID] = basketItem
	}

	basketItem.Count += input.Count

	_, basketRepositorySaveErr := useCase.basketRepository.Save(userBasket)
	if basketRepositorySaveErr != nil {
		return nil, basketRepositorySaveErr
	}

	userBasketDTO, basketOutputServiceErr := useCase.basketOutputService.CreateBasketDTO(userBasket)
	if basketOutputServiceErr != nil {
		return nil, basketOutputServiceErr
	}

	output := &AddProductUseCaseOutput{
		UserBasketDTO: userBasketDTO,
		Actions:       map[string]string{},
	}

	return output, nil
}
