package inmemory

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/entities"
)

var _ entities.BasketRepository = (*MongoBasketRepository)(nil)

type MongoBasketRepository struct {
	collection *mongo.Collection
}

func NewMongoBasketRepository(collection *mongo.Collection) entities.BasketRepository {
	return &MongoBasketRepository{
		collection: collection,
	}
}

func (repository *MongoBasketRepository) Save(basket *entities.Basket) (string, error) {
	if basket == nil {
		return "", fmt.Errorf("basket is nil")
	}

	if basket.ID == "" {
		basket.ID = uuid.NewString()
	}

	result, replaceErr := repository.collection.ReplaceOne(context.Background(), bson.M{"id": basket.ID}, basket)
	if replaceErr != nil {
		return "", replaceErr
	}

	if result.MatchedCount == 0 {
		_, err := repository.collection.InsertOne(context.Background(), basket)
		if err != nil {
			return "", err
		}
	}

	return basket.ID, nil
}

func (repository *MongoBasketRepository) Find(id string) (*entities.Basket, error) {
	result := repository.collection.FindOne(context.Background(), bson.M{"id": id})
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return nil, &entities.BasketNotFoundError{}
		}
		return nil, result.Err()
	}

	var basket entities.Basket
	decodeErr := result.Decode(&basket)
	if decodeErr != nil {
		return nil, decodeErr
	}

	return &basket, nil
}

func (repository *MongoBasketRepository) FindByUserId(userId string) (*entities.Basket, error) {
	result := repository.collection.FindOne(context.Background(), bson.M{"userid": userId})
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return nil, &entities.BasketNotFoundError{}
		}

		return nil, result.Err()
	}

	var basket entities.Basket
	decodeErr := result.Decode(&basket)
	if decodeErr != nil {
		return nil, decodeErr
	}

	return &basket, nil
}
