package inmemory

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/entities"
)

func initTestcontainers(t *testing.T) (string, func()) {
	ctx := context.Background()

	mongodbContainer, err := mongodb.Run(ctx, "mongodb/mongodb-community-server:8.0-ubi8")
	stop := func() {
		if err := testcontainers.TerminateContainer(mongodbContainer); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}

	require.NoError(t, err)
	require.NotNil(t, mongodbContainer)

	endpoint, err := mongodbContainer.ConnectionString(ctx)
	if err != nil {
		log.Fatalf("failed to get connection string: %s", err)
	}

	return endpoint, stop
}
func Test_MongoBasketRepository(t *testing.T) {
	endpoint, stop := initTestcontainers(t)
	defer stop()

	clientOpts := options.Client().ApplyURI(endpoint)
	mongoClient, mongoClientErr := mongo.Connect(clientOpts)
	if mongoClientErr != nil {
		panic(mongoClientErr)
	}
	defer func() {
		if mongoClientErr = mongoClient.Disconnect(context.TODO()); mongoClientErr != nil {
			panic(mongoClientErr)
		}
	}()

	basketsCollection := mongoClient.Database(DatabaseName).Collection(BasketsCollectionName)
	repository := NewMongoBasketRepository(basketsCollection)

	foundBasket1, err := repository.Find("1")

	require.Error(t, err)
	require.Nil(t, foundBasket1)

	basketId := "1"
	userID := "1337"

	basketFactory := entities.NewBasketFactory()
	basket, basketErr := basketFactory.NewBasketWithID(basketId, userID)

	require.NoError(t, basketErr)
	require.NotNil(t, basket)

	returnedBasketId, err := repository.Save(basket)

	require.NoError(t, err)
	require.NotEmpty(t, returnedBasketId)

	foundBasket, err := repository.Find(basketId)

	require.NoError(t, err)
	require.NotNil(t, foundBasket)
	require.Equal(t, basket, foundBasket)
	require.Len(t, foundBasket.GetItems(), 0)

	_ = basket.AddItem("A12345", 1)

	basketId, err = repository.Save(basket)

	require.NoError(t, err)
	require.NotEmpty(t, basketId)

	foundBasket2, err := repository.FindByUserId("1336")

	require.Error(t, err)
	require.Nil(t, foundBasket2)

	foundBasket3, err := repository.FindByUserId(userID)

	require.NoError(t, err)
	require.NotNil(t, foundBasket3)
	require.Equal(t, basket, foundBasket3)
	require.Len(t, foundBasket3.GetItems(), 1)
}
