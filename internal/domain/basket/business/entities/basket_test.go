package entities

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Basket(t *testing.T) {
	basketID := "1"
	userID := "1337"

	factory := NewBasketFactory()
	basket, err := factory.NewBasketWithID(basketID, userID)

	require.NoError(t, err)
	require.NotNil(t, basket)

	require.Equal(t, basketID, basket.GetID())
	require.Equal(t, userID, basket.GetUserID())

	require.NotNil(t, basket.GetItems())
	require.Empty(t, basket.GetItems())

}

func Test_Basket_WithItems(t *testing.T) {
	basketID := "1"
	userID := "1337"

	factory := NewBasketFactory()
	basket, err := factory.NewBasketWithID(basketID, userID)

	require.NoError(t, err)
	require.NotNil(t, basket)

	require.Equal(t, basketID, basket.GetID())
	require.Equal(t, userID, basket.GetUserID())

	require.NotNil(t, basket.GetItems())
	require.Empty(t, basket.GetItems())

	productID := "A12345"
	productCount := 5

	require.False(t, basket.HasItem(productID))

	basketItem, err := basket.GetItem(productID)
	require.Error(t, err)
	require.Nil(t, basketItem)

	basket.AddItem(productID, productCount)
	require.Equal(t, 1, len(basket.GetItems()))

	require.True(t, basket.HasItem(productID))

	basketItem, err = basket.GetItem(productID)
	require.NoError(t, err)
	require.NotNil(t, basketItem)

	require.Equal(t, productID, basketItem.GetProductID())
	require.Equal(t, productCount, basketItem.GetCount())

	basketItem.SetCount(productCount + 1)
	require.Equal(t, productCount+1, basketItem.GetCount())

	_ = basket.RemoveItem(productID)
	require.Equal(t, 0, len(basket.GetItems()))

	require.False(t, basket.HasItem(productID))
}

func Test_Basket_ClearBasket(t *testing.T) {
	basketID := "1"
	userID := "1337"

	factory := NewBasketFactory()
	basket, err := factory.NewBasketWithID(basketID, userID)

	require.NoError(t, err)
	require.NotNil(t, basket)

	require.Equal(t, basketID, basket.GetID())
	require.Equal(t, userID, basket.GetUserID())

	require.NotNil(t, basket.GetItems())
	require.Empty(t, basket.GetItems())

	productID := "A12345"
	productCount := 5

	basket.AddItem(productID, productCount)
	require.Equal(t, 1, len(basket.GetItems()))

	basket.Clear()
	require.Equal(t, 0, len(basket.GetItems()))
}
