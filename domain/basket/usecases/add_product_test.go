package usecases

import (
	"github.com/arkadiusjonczek/clean-architecture-go/domain/basket/usecases/helper"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/arkadiusjonczek/clean-architecture-go/domain/basket"
	"github.com/arkadiusjonczek/clean-architecture-go/domain/warehouse"
)

func Test_AddProductToBasketUseCase_WrongInput_ReturnsError(t *testing.T) {
	testCases := map[string]struct {
		input *AddProductUseCaseInput
	}{
		"input is nil": {
			input: nil,
		},
		"UserID is empty": {
			input: &AddProductUseCaseInput{},
		},
		"ProductID is empty": {
			input: &AddProductUseCaseInput{
				UserID: "1337",
			},
		},
		"Count is invalid": {
			input: &AddProductUseCaseInput{
				UserID:    "1337",
				ProductID: "1",
			},
		},
	}

	for errorString, testCase := range testCases {
		t.Run(errorString, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			basketFactory := basket.NewBasketFactory()
			basketRepositoryMock := basket.NewMockBasketRepository(ctrl)
			productRepositoryMock := warehouse.NewMockProductRepository(ctrl)

			basketService := helper.NewBasketCreatorServiceImpl(basketFactory, basketRepositoryMock)

			useCase := NewAddProductUseCaseImpl(basketService, basketRepositoryMock, productRepositoryMock)

			_, err := useCase.Execute(testCase.input)

			require.Error(t, err)
			require.ErrorContains(t, err, errorString)
		})
	}
}

// TODO: Test also use cases like product not found, product out of stock etc.
func Test_AddProductToBasketUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	basketID := "12345"
	userID := "1337"

	productID := "1"
	productName := "Product 1"
	product1 := &warehouse.Product{
		ID:    productID,
		Name:  productName,
		Stock: 10,
		Price: &warehouse.ProductPrice{
			Price:    13.37,
			Currency: "EUR",
		},
	}

	basketFactory := basket.NewBasketFactory()

	basketRepositoryMock := basket.NewMockBasketRepository(ctrl)

	userBasket, err := basketFactory.NewBasket(userID)
	userBasket.ID = basketID // set basket id, otherwise a new basket will be created
	require.NoError(t, err)

	basketRepositoryMock.EXPECT().FindByUserId(userID).Return(userBasket, nil)
	basketRepositoryMock.EXPECT().Save(&basket.Basket{
		ID:     basketID,
		UserID: userID,
		Items: map[string]*basket.BasketItem{
			productID: {
				Product: product1,
				Count:   1,
			},
		},
	}).Return(basketID, nil)

	productRepositoryMock := warehouse.NewMockProductRepository(ctrl)
	productRepositoryMock.EXPECT().Find(productID).Return(product1, nil)

	basketService := helper.NewBasketCreatorServiceImpl(basketFactory, basketRepositoryMock)

	useCase := NewAddProductUseCaseImpl(basketService, basketRepositoryMock, productRepositoryMock)

	input := &AddProductUseCaseInput{
		UserID:    userID,
		ProductID: productID,
		Count:     1,
	}

	output, err := useCase.Execute(input)

	require.NoError(t, err)
	require.NotNil(t, output)
}
