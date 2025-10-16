package helper

//go:generate mockgen -source=product_price_simulator_service.go -destination=product_price_simulator_service_mock.go -package=helper

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/warehouse/business/entities"
)

// ProductPriceSimulatorService will make price changes to demonstrate the difference between Basket and BasketDTO
type ProductPriceSimulatorService interface {
	Execute()
}

var _ ProductPriceSimulatorService = (*ProductPriceSimulatorServiceImpl)(nil)

type ProductPriceSimulatorServiceImpl struct {
	productRepository entities.ProductRepository
}

func NewProductPriceSimulatorService(productRepository entities.ProductRepository) (ProductPriceSimulatorService, error) {
	if productRepository == nil {
		return nil, fmt.Errorf("productRepository is nil")
	}

	return &ProductPriceSimulatorServiceImpl{
		productRepository: productRepository,
	}, nil
}

func (service *ProductPriceSimulatorServiceImpl) Execute() {
	plus := rand.Intn(2) == 0
	for _, product := range service.productRepository.FindAll() {
		change := rand.Float64() * 0.10
		oldPrice := product.Price.Value
		if plus {
			product.Price.Value += change
		} else {
			product.Price.Value -= change
		}
		log.Printf("ProductPriceSimulatorService: Updating Product %s price: %f (old price: %f)\n", product.ID, product.Price.Value, oldPrice)
		service.productRepository.Save(product)
	}
}
