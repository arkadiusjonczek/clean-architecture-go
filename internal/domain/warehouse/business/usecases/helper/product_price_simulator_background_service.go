package helper

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

const (
	ProductPriceSimulatorWaitDuration = 10 * time.Second
)

// ProductPriceSimulatorBackgroundService will make price changes to demonstrate the difference between Basket and BasketDTO
type ProductPriceSimulatorBackgroundService interface {
	Start()
	Stop()
}

var _ ProductPriceSimulatorBackgroundService = (*ProductPriceSimulatorBackgroundServiceImpl)(nil)

type ProductPriceSimulatorBackgroundServiceImpl struct {
	ctx                          context.Context
	started                      bool
	syncMutex                    sync.Mutex
	productPriceSimulatorService ProductPriceSimulatorService
}

func NewProductPriceSimulatorBackgroundService(productPriceSimulatorService ProductPriceSimulatorService) (ProductPriceSimulatorBackgroundService, error) {
	if productPriceSimulatorService == nil {
		return nil, fmt.Errorf("productPriceSimulatorService is nil")
	}

	return &ProductPriceSimulatorBackgroundServiceImpl{
		ctx:                          context.Background(),
		started:                      false,
		productPriceSimulatorService: productPriceSimulatorService,
	}, nil
}

// Start is not thread safe
func (service *ProductPriceSimulatorBackgroundServiceImpl) Start() {
	service.syncMutex.Lock()
	defer service.syncMutex.Unlock()

	if service.started {
		return
	}

	log.Println("ProductPriceSimulatorBackgroundService: Starting...")
	go service.start()
	service.started = true
}

func (service *ProductPriceSimulatorBackgroundServiceImpl) start() {
outerloop:
	for {
		service.productPriceSimulatorService.Execute()

		select {
		case <-service.ctx.Done():
			break outerloop
		case <-time.After(ProductPriceSimulatorWaitDuration):
		}
	}
}

func (service *ProductPriceSimulatorBackgroundServiceImpl) Stop() {
	service.syncMutex.Lock()
	defer service.syncMutex.Unlock()

	if !service.started {
		return
	}

	log.Println("ProductPriceSimulatorBackgroundService: Stopping...")
	service.ctx.Done()
}
