package inmemory

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/entities"
)

func Test_InMemoryBasketRepository_Find(t *testing.T) {
	repository := NewInMemoryBasketRepository()

	require.NotNil(t, repository)

	basketID := "1"
	userID := "1337"

	basket, err := repository.Find(basketID)

	require.Error(t, err)
	require.Nil(t, basket)

	factory := entities.NewBasketFactory()

	basket, basketErr := factory.NewBasketWithID(basketID, userID)

	require.NoError(t, basketErr)

	repositoryBasketID, repositoryErr := repository.Save(basket)

	require.Nil(t, repositoryErr)
	require.Equal(t, basketID, repositoryBasketID)

	basket, err = repository.Find(basketID)

	require.NoError(t, err)
	require.NotNil(t, basket)
	require.Equal(t, basketID, basket.GetID())
	require.Equal(t, userID, basket.GetUserID())
	require.NotNil(t, basket.GetItems())
	require.Empty(t, basket.GetItems())
}

func Test_InMemoryBasketRepository_FindByUserId(t *testing.T) {
	repository := NewInMemoryBasketRepository()

	require.NotNil(t, repository)

	basketID := "1"
	userID := "1337"

	basket, err := repository.FindByUserId(userID)

	require.Error(t, err)
	require.Nil(t, basket)

	factory := entities.NewBasketFactory()

	basket, basketErr := factory.NewBasketWithID(basketID, userID)

	require.NoError(t, basketErr)

	repositoryBasketID, repositoryErr := repository.Save(basket)

	require.Nil(t, repositoryErr)
	require.Equal(t, basketID, repositoryBasketID)

	basket, err = repository.FindByUserId(userID)

	require.NoError(t, err)
	require.NotNil(t, basket)
	require.Equal(t, basketID, basket.GetID())
	require.Equal(t, userID, basket.GetUserID())
	require.NotNil(t, basket.GetItems())
	require.Empty(t, basket.GetItems())
}

func Test_InMemoryBasketRepository_SaveNewBasket(t *testing.T) {
	repository := NewInMemoryBasketRepository()

	require.NotNil(t, repository)

	userID := "1337"

	basket, err := repository.FindByUserId(userID)

	require.Error(t, err)
	require.Nil(t, basket)

	factory := entities.NewBasketFactory()

	basket, basketErr := factory.NewBasket(userID)

	require.NoError(t, basketErr)

	repositoryBasketID, repositoryErr := repository.Save(basket)

	require.Nil(t, repositoryErr)
	require.NotEmpty(t, repositoryBasketID)

	basket, err = repository.FindByUserId(userID)

	require.NoError(t, err)
	require.NotNil(t, basket)
	require.Equal(t, repositoryBasketID, basket.GetID())
	require.Equal(t, userID, basket.GetUserID())
	require.NotNil(t, basket.GetItems())
	require.Empty(t, basket.GetItems())
}
