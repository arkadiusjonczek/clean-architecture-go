package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/adapters/rest"
	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/adapters/web"
	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/entities"
	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/usecases"
	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/usecases/helper"
	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/drivers/inmemory"
	basketdrivermongodb "github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/drivers/mongodb"
	basketdrivermysql "github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/drivers/mysql"
	warehouse "github.com/arkadiusjonczek/clean-architecture-go/internal/domain/warehouse/business/entities"
	warehousehelper "github.com/arkadiusjonczek/clean-architecture-go/internal/domain/warehouse/business/usecases/helper"
	warehousedriverinmemory "github.com/arkadiusjonczek/clean-architecture-go/internal/domain/warehouse/drivers/inmemory"
)

func main() {
	err := startHTTPServer()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func startHTTPServer() error {
	fmt.Println("Starting server...")

	// create drivers

	var basketRepository entities.BasketRepository

	switch os.Getenv("DRIVER") {
	case "mongodb":
		fmt.Printf("Driver: MongoDB\n")

		clientOpts := options.Client().ApplyURI("mongodb://localhost:27017")
		mongoClient, mongoClientErr := mongo.Connect(clientOpts)
		if mongoClientErr != nil {
			panic(mongoClientErr)
		}
		defer func() {
			if mongoClientErr = mongoClient.Disconnect(context.TODO()); mongoClientErr != nil {
				panic(mongoClientErr)
			}
		}()

		basketsCollection := mongoClient.Database("ecommerce").Collection("baskets")

		basketRepository = basketdrivermongodb.NewMongoBasketRepository(basketsCollection)
	case "mysql":
		fmt.Printf("Driver: MySQL\n")

		// Get MySQL configuration from environment variables
		mysqlHost := os.Getenv("MYSQL_HOST")
		if mysqlHost == "" {
			mysqlHost = "localhost"
		}
		mysqlPort := os.Getenv("MYSQL_PORT")
		if mysqlPort == "" {
			mysqlPort = "3306"
		}
		mysqlUsername := os.Getenv("MYSQL_USERNAME")
		if mysqlUsername == "" {
			mysqlUsername = "root"
		}
		mysqlPassword := os.Getenv("MYSQL_PASSWORD")
		if mysqlPassword == "" {
			mysqlPassword = "password"
		}
		mysqlDatabase := os.Getenv("MYSQL_DATABASE")
		if mysqlDatabase == "" {
			mysqlDatabase = "ecommerce"
		}

		// Create MySQL database connection
		config := basketdrivermysql.DatabaseConfig{
			Host:     mysqlHost,
			Port:     mysqlPort,
			Username: mysqlUsername,
			Password: mysqlPassword,
			Database: mysqlDatabase,
		}

		mysqlDB, mysqlErr := basketdrivermysql.NewConnection(config)
		if mysqlErr != nil {
			panic(fmt.Errorf("failed to connect to MySQL: %w", mysqlErr))
		}
		defer func() {
			if closeErr := mysqlDB.Close(); closeErr != nil {
				log.Printf("Failed to close MySQL connection: %v", closeErr)
			}
		}()

		basketRepository = basketdrivermysql.NewMySQLBasketRepository(mysqlDB)
	default:
		fmt.Printf("Driver: InMemory\n")

		basketRepository = inmemory.NewInMemoryBasketRepository()
	}

	productRepository := warehousedriverinmemory.NewInMemoryProductRepository()
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
			Stock: 30,
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
			Stock: 0,
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
			Stock: 50,
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

	// simulate price changes

	productPriceSimulatorService, productPriceSimulatorServiceErr := warehousehelper.NewProductPriceSimulatorService(productRepository)
	if productPriceSimulatorServiceErr != nil {
		return productPriceSimulatorServiceErr
	}

	productPriceSimulatorBackgroundService, productPriceSimulatorBackgroundServiceErr := warehousehelper.NewProductPriceSimulatorBackgroundService(productPriceSimulatorService)
	if productPriceSimulatorBackgroundServiceErr != nil {
		return productPriceSimulatorBackgroundServiceErr
	}

	productPriceSimulatorBackgroundService.Start()
	defer productPriceSimulatorBackgroundService.Stop()

	// create interface adapters

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	webBasketController := web.NewBasketController(showBasketUseCase)
	webBasketControllerRouter := web.NewBasketControllerRouter(webBasketController)
	webBasketControllerRouterErr := webBasketControllerRouter.RegisterRoutes(router)
	if webBasketControllerRouterErr != nil {
		return webBasketControllerRouterErr
	}

	restBasketController := rest.NewBasketController(showBasketUseCase, clearBasketUseCase, addProductUseCase, updateProductCountUseCase, removeProductUseCase)
	restBasketControllerRouter := rest.NewBasketControllerRouter(restBasketController)
	restBasketControllerRouterErr := restBasketControllerRouter.RegisterRoutes(router)
	if restBasketControllerRouterErr != nil {
		return restBasketControllerRouterErr
	}

	// start http server

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = "localhost:8080"
	}

	fmt.Printf("Starting server on %s\n", addr)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to listen and serve: %v", err)
		}
	}()
	defer func(httpServer *http.Server) {
		err := httpServer.Close()
		if err != nil {
			log.Fatalf("Failed to close server: %v", err)
		}
	}(httpServer)

	// Wait for SIGINT (Ctrl+C) signal to exit
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	fmt.Println("Stopping server...")

	return nil
}
