package usecases

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/arkadiusjonczek/clean-architecture-go/domain/basket"
	"github.com/arkadiusjonczek/clean-architecture-go/domain/warehouse"
)

// TODO: Test also use cases like product not found, product out of stock etc.
func Test_AddProductToBasketUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	basketID := "12345"
	userID := "1337"

	productID := "1"
	productName := "Product 1"

	basketRepositoryMock := basket.NewMockBasketRepository(ctrl)

	userBasket, err := basket.NewBasket(basketID, userID)
	require.NoError(t, err)

	basketRepositoryMock.EXPECT().FindByUserId(userID).Return(userBasket, nil)
	basketRepositoryMock.EXPECT().Save(&basket.Basket{
		ID:     basketID,
		UserID: userID,
		Items: map[string]*basket.BasketItem{
			productID: {
				Product: warehouse.Product{
					ID:   productID,
					Name: productName,
				},
				Count: 1,
			},
		},
	})

	productRepositoryMock := warehouse.NewMockProductRepository(ctrl)
	productRepositoryMock.EXPECT().Find(productID).Return(&warehouse.Product{
		ID:   productID,
		Name: productName,
	}, nil)

	useCase := NewAddProductToBasketUseCaseImpl(basketRepositoryMock, productRepositoryMock)

	input := &AddProductToBasketUseCaseInput{
		UserID:    userID,
		ProductID: productID,
		Count:     1,
	}

	output, err := useCase.Execute(input)

	require.NoError(t, err)
	require.NotNil(t, output)
}
