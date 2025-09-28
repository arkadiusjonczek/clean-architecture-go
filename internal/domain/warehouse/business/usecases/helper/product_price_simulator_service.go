package helper

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/warehouse/business/entities"
)

const (
	PRODUCT_PRICE_SIMULATOR_WAIT_DURATION = 10 * time.Second
)

// ProductPriceSimulatorService will make price changes to demonstrate the difference between Basket and BasketDTO
type ProductPriceSimulatorService interface {
	Start()
	Stop()
}

var _ ProductPriceSimulatorService = (*ProductPriceSimulatorServiceImpl)(nil)

type ProductPriceSimulatorServiceImpl struct {
	ctx               context.Context
	productRepository entities.ProductRepository
}

func NewProductPriceSimulator(productRepository entities.ProductRepository) ProductPriceSimulatorService {
	return &ProductPriceSimulatorServiceImpl{
		ctx:               context.Background(),
		productRepository: productRepository,
	}
}

// Start is not thread safe
func (service *ProductPriceSimulatorServiceImpl) Start() {
	log.Println("ProductPriceSimulatorService: Starting...")
	go service.start()
}

func (service *ProductPriceSimulatorServiceImpl) start() {
	for {
		service.execute()

		select {
		case <-service.ctx.Done():
			break
		case <-time.After(PRODUCT_PRICE_SIMULATOR_WAIT_DURATION):
		}
	}
}

func (service *ProductPriceSimulatorServiceImpl) execute() {
	plus := rand.Intn(1) == 0
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

func (service *ProductPriceSimulatorServiceImpl) Stop() {
	log.Println("ProductPriceSimulatorService: Stopping...")
	service.ctx.Done()
}
