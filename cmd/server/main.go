package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/entities"
	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/usecases"
	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/usecases/helper"
	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/controllers/rest"
	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/drivers/inmemory"
	warehouse "github.com/arkadiusjonczek/clean-architecture-go/internal/domain/warehouse/business/entities"
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
			ID:   "1",
			Name: "Product 1",
			Price: &warehouse.ProductPrice{
				Price:    11.99,
				Currency: "EUR",
			},
			Stock: 10,
		},
	)
	productRepository.Save(
		&warehouse.Product{
			ID:   "2",
			Name: "Product 2",
			Price: &warehouse.ProductPrice{
				Price:    12.99,
				Currency: "EUR",
			},
			Stock: 20,
		},
	)

	// create business logic and inject drivers

	basketFactory := entities.NewBasketFactory()

	basketService := helper.NewBasketCreatorServiceImpl(basketFactory, basketRepository)

	showBasketUseCase := usecases.NewShowBasketUseCaseImpl(basketService)
	clearBasketUseCase := usecases.NewClearBasketUseCaseImpl(basketService, basketRepository)
	addProductUseCase := usecases.NewAddProductUseCaseImpl(basketService, basketRepository, productRepository)
	updateProductCountUseCase := usecases.NewUpdateProductCountImpl(basketService, basketRepository, productRepository)
	removeProductUseCase := usecases.NewRemoveProductUseCaseImpl(basketService, basketRepository, productRepository)

	// create interface adapters

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	basketController := rest.NewBasketController(showBasketUseCase, clearBasketUseCase, addProductUseCase, updateProductCountUseCase, removeProductUseCase)

	basketControllerRouter := rest.NewBasketControllerRouter(basketController)
	basketControllerRouter.RegisterRoutes(router)

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
