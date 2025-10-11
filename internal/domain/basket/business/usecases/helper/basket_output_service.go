package helper

import (
	"fmt"

	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/entities"
	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/usecases/dto"
	warehouse "github.com/arkadiusjonczek/clean-architecture-go/internal/domain/warehouse/business/entities"
)

type BasketOutputService interface {
	CreateBasketDTO(basket *entities.Basket) (*dto.BasketDTO, error)
}

var _ BasketOutputService = (*BasketOutputServiceImpl)(nil)

type BasketOutputServiceImpl struct {
	productRepository warehouse.ProductRepository
}

func NewBasketOutputService(productRepository warehouse.ProductRepository) BasketOutputService {
	return &BasketOutputServiceImpl{
		productRepository: productRepository,
	}
}

func (service *BasketOutputServiceImpl) CreateBasketDTO(basket *entities.Basket) (*dto.BasketDTO, error) {
	if basket == nil {
		return nil, fmt.Errorf("basket is nil")
	}

	basketDTO := &dto.BasketDTO{
		Items: []*dto.BasketItem{},
	}

	// order guarantee
	basketItemsKeys := make([]string, 0)
	for k, _ := range basket.GetItems() {
		basketItemsKeys = append(basketItemsKeys, k)
	}
	//sort.Strings(basketItemsKeys)

	for _, productId := range basketItemsKeys {
		item, _ := basket.GetItem(productId)

		product, productRepositoryErr := service.productRepository.Find(item.ProductID)
		if productRepositoryErr != nil {
			return nil, productRepositoryErr
		}

		basketProduct := &dto.Product{
			ID:   product.ID,
			Name: product.Name,
			Price: &dto.ProductPrice{
				Value:    fmt.Sprintf("%.2f", product.Price.Value),
				Currency: product.Price.Currency,
			},
		}

		basketItem := &dto.BasketItem{
			Product: basketProduct,
			Count:   item.Count,
		}

		basketDTO.Items = append(basketDTO.Items, basketItem)
	}

	return basketDTO, nil
}
