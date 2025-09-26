package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/adapters/rest"
	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/adapters/web"
	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/entities"
	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/usecases"
	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/usecases/helper"
	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/drivers/inmemory"
	warehouse "github.com/arkadiusjonczek/clean-architecture-go/internal/domain/warehouse/business/entities"
	warehousehelper "github.com/arkadiusjonczek/clean-architecture-go/internal/domain/warehouse/business/usecases/helper"
)

func main() {
	startHTTPServer()
}

func startHTTPServer() {
	fmt.Println("Starting server...")

	// create drivers

	basketRepository := inmemory.NewInMemoryBasketRepository()

	productRepository := inmemory.NewInMemoryProductRepository()
	productRepository.Save(
		&warehouse.Product{
			ID:   "A12341",
			Name: "Product 1",
			Price: &warehouse.ProductPrice{
				Value:    11.99,
				Currency: "EUR",
			},
			Stock: 10,
		},
	)
	productRepository.Save(
		&warehouse.Product{
			ID:   "A12342",
			Name: "Product 2",
			Price: &warehouse.ProductPrice{
				Value:    12.99,
				Currency: "EUR",
			},
			Stock: 20,
		},
	)
	productRepository.Save(
		&warehouse.Product{
			ID:   "A12343",
			Name: "Product 3",
			Price: &warehouse.ProductPrice{
				Value:    13.99,
				Currency: "EUR",
			},
			Stock: 0,
		},
	)
	productRepository.Save(
		&warehouse.Product{
			ID:   "A12344",
			Name: "Product 4",
			Price: &warehouse.ProductPrice{
				Value:    14.99,
				Currency: "EUR",
			},
			Stock: 40,
		},
	)
	productRepository.Save(
		&warehouse.Product{
			ID:   "A12345",
			Name: "Product 5",
			Price: &warehouse.ProductPrice{
				Value:    15.99,
				Currency: "EUR",
			},
			Stock: 0,
		},
	)

	// create business logic and inject drivers

	basketFactory := entities.NewBasketFactory()

	basketCreatorService := helper.NewBasketCreatorServiceImpl(basketFactory, basketRepository)
	basketOutputService := helper.NewBasketOutputService(productRepository)

	showBasketUseCase := usecases.NewShowBasketUseCaseImpl(basketCreatorService, basketOutputService)
	clearBasketUseCase := usecases.NewClearBasketUseCaseImpl(basketCreatorService, basketOutputService, basketRepository)
	addProductUseCase := usecases.NewAddProductUseCaseImpl(basketCreatorService, basketOutputService, basketRepository, productRepository)
	updateProductCountUseCase := usecases.NewUpdateProductCountImpl(basketCreatorService, basketOutputService, basketRepository, productRepository)
	removeProductUseCase := usecases.NewRemoveProductUseCaseImpl(basketCreatorService, basketOutputService, basketRepository, productRepository)

	productPriceSimulatorService := warehousehelper.NewProductPriceSimulator(productRepository)
	productPriceSimulatorService.Start()
	defer productPriceSimulatorService.Stop()

	// create interface adapters

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	webBasketController := web.NewBasketController(showBasketUseCase)
	webBasketControllerRouter := web.NewBasketControllerRouter(webBasketController)
	webBasketControllerRouter.RegisterRoutes(router)

	restBasketController := rest.NewBasketController(showBasketUseCase, clearBasketUseCase, addProductUseCase, updateProductCountUseCase, removeProductUseCase)
	restBasketControllerRouter := rest.NewBasketControllerRouter(restBasketController)
	restBasketControllerRouter.RegisterRoutes(router)

	// start http server

	httpServer := &http.Server{
		Addr:    "localhost:8080",
		Handler: router,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to listen and serve: %v", err)
		}
	}()
	defer httpServer.Close()

	// Wait for SIGINT (Ctrl+C) signal to exit
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	fmt.Println("Stopping server...")
}
