package usecases

import (
	"fmt"

	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/entities"
	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/usecases/dto"
	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/usecases/helper"
	warehouse "github.com/arkadiusjonczek/clean-architecture-go/internal/domain/warehouse/business/entities"
)

type UpdateProductCountUseCaseInput struct {
	UserID    string
	ProductID string
	Count     int
}

type UpdateProductCountUseCaseOutput struct {
	UserBasket *dto.BasketDTO
	Actions    map[string]string
}

type UpdateProductCountUseCase interface {
	Execute(input *UpdateProductCountUseCaseInput) (*UpdateProductCountUseCaseOutput, error)
}

func NewUpdateProductCountImpl(basketService helper.BasketCreatorService, basketOutputService helper.BasketOutputService, basketRepository entities.BasketRepository, productRepository warehouse.ProductRepository) UpdateProductCountUseCase {
	return &UpdateProductCountUseCaseImpl{
		basketService:       basketService,
		basketOutputService: basketOutputService,
		basketRepository:    basketRepository,
		productRepository:   productRepository,
	}
}

var _ UpdateProductCountUseCase = (*UpdateProductCountUseCaseImpl)(nil)

type UpdateProductCountUseCaseImpl struct {
	basketService       helper.BasketCreatorService
	basketOutputService helper.BasketOutputService
	basketRepository    entities.BasketRepository
	productRepository   warehouse.ProductRepository
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

	if product.Stock <= 0 {
		return nil, fmt.Errorf("product %s is out of stock", input.ProductID)
	}

	var basketItem *entities.BasketItem
	if !userBasket.HasItem(input.ProductID) {
		// add the product with count = 0
		// because of the following product stock check
		// which will set the correct count
		// related to the available stock
		basketItem = userBasket.AddItem(input.ProductID, 0)
	} else {
		basketItem, _ = userBasket.GetItem(input.ProductID)
	}

	var actions map[string]string
	if product.Stock < input.Count {
		basketItem.Count = product.Stock
		actions = map[string]string{
			"product_stock": fmt.Sprintf("Product %s stock is too low to add %d. Updated basket item count to %d.", input.ProductID, input.Count, product.Stock),
		}
	} else {
		basketItem.Count = input.Count
	}

	_, basketRepositorySaveErr := useCase.basketRepository.Save(userBasket)
	if basketRepositorySaveErr != nil {
		return nil, basketRepositorySaveErr
	}

	userBasketDTO, basketOutputServiceErr := useCase.basketOutputService.CreateBasketDTO(userBasket)
	if basketOutputServiceErr != nil {
		return nil, basketOutputServiceErr
	}

	output := &UpdateProductCountUseCaseOutput{
		UserBasket: userBasketDTO,
		Actions:    actions,
	}

	return output, nil
}
