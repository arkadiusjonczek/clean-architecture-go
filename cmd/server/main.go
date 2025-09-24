package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/arkadiusjonczek/clean-architecture-go/domain/basket"
	"github.com/arkadiusjonczek/clean-architecture-go/domain/basket/usecases"
	"github.com/arkadiusjonczek/clean-architecture-go/domain/basket/usecases/helper"
	"github.com/arkadiusjonczek/clean-architecture-go/domain/warehouse"
	"github.com/arkadiusjonczek/clean-architecture-go/internal/adapters/rest"
	"github.com/arkadiusjonczek/clean-architecture-go/internal/drivers/inmemory"
)

func main() {
	startHTTPServer()
}

func startHTTPServer() {
	fmt.Println("Starting server...")

	//r := mux.NewRouter()
	//r.HandleFunc("/", HomeHandler)
	//http.Handle("/", r)

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// create dependencies

	basketFactory := basket.NewBasketFactory()
	basketRepository := inmemory.NewInMemoryBasketRepository()

	basketService := helper.NewBasketCreatorServiceImpl(basketFactory, basketRepository)

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

	showBasketUseCase := usecases.NewShowBasketUseCaseImpl(basketService)
	clearBasketUseCase := usecases.NewClearBasketUseCaseImpl(basketService, basketRepository)
	addProductUseCase := usecases.NewAddProductUseCaseImpl(basketService, basketRepository, productRepository)
	updateProductCountUseCase := usecases.NewUpdateProductCountImpl(basketService, basketRepository, productRepository)
	removeProductUseCase := usecases.NewRemoveProductUseCaseImpl(basketService, basketRepository, productRepository)

	basketController := rest.NewBasketController(showBasketUseCase, clearBasketUseCase, addProductUseCase, updateProductCountUseCase, removeProductUseCase)

	basketControllerRouter := rest.NewBasketControllerRouter(basketController)
	basketControllerRouter.RegisterRoutes(router)

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
