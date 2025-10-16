package entities

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_BasketFactory_NewBasket_ReturnsError(t *testing.T) {
	factory := NewBasketFactory()
	basket, err := factory.NewBasket("")

	require.Error(t, err)
	require.ErrorContains(t, err, "userID cannot be empty")
	require.Nil(t, basket)
}

func Test_BasketFactory_NewBasket(t *testing.T) {
	userID := "1337"

	factory := NewBasketFactory()
	basket, err := factory.NewBasket(userID)

	require.NoError(t, err)
	require.NotNil(t, basket)

	require.Equal(t, "", basket.GetID())
	require.Equal(t, userID, basket.GetUserID())
}

func Test_BasketFactory_NewBasketWithID_ReturnsError(t *testing.T) {
	testCases := map[string]struct {
		basketID string
		userID   string
	}{
		"basketID cannot be empty": {
			basketID: "",
			userID:   "",
		},
		"userID cannot be empty": {
			basketID: "1",
			userID:   "",
		},
	}

	for errorMessage, testCase := range testCases {
		t.Run("", func(t *testing.T) {
			factory := NewBasketFactory()
			basket, err := factory.NewBasketWithID(testCase.basketID, testCase.userID)

			require.Error(t, err)
			require.ErrorContains(t, err, errorMessage)
			require.Nil(t, basket)
		})
	}
}

func Test_BasketFactory_NewBasketWithID(t *testing.T) {
	basketID := "1"
	userID := "1337"

	factory := NewBasketFactory()
	basket, err := factory.NewBasketWithID(basketID, userID)

	require.NoError(t, err)
	require.NotNil(t, basket)

	require.Equal(t, basketID, basket.GetID())
	require.Equal(t, userID, basket.GetUserID())
}
