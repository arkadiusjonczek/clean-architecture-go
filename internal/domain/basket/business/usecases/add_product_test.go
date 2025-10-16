package usecases

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/entities"
	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/usecases/helper"
	warehouse "github.com/arkadiusjonczek/clean-architecture-go/internal/domain/warehouse/business/entities"
)

func Test_AddProductToBasketUseCase_NewAddProductUseCaseImpl_ReturnsError(t *testing.T) {
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

			basketFactory := entities.NewBasketFactory()
			basketRepositoryMock := entities.NewMockBasketRepository(ctrl)
			productRepositoryMock := warehouse.NewMockProductRepository(ctrl)

			basketCreatorService := helper.NewBasketCreatorServiceImpl(basketFactory, basketRepositoryMock)
			basketOutputService := helper.NewBasketOutputService(productRepositoryMock)

			useCase := NewAddProductUseCaseImpl(basketCreatorService, basketOutputService, basketRepositoryMock, productRepositoryMock)

			_, err := useCase.Execute(testCase.input)

			require.Error(t, err)
			require.ErrorContains(t, err, errorString)
		})
	}
}

// TODO: Test also use cases like product not found, product out of stock etc.
func Test_AddProductToBasketUseCase(t *testing.T) {
	// arrange

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	basketID := "12345"
	userID := "1337"

	product1ID := "1"
	product1Name := "Product 1"
	product1 := &warehouse.Product{
		ID:    product1ID,
		Name:  product1Name,
		Stock: 10,
		Price: &warehouse.ProductPrice{
			Value:    13.37,
			Currency: "EUR",
		},
	}

	basketFactory := entities.NewBasketFactory()

	basketRepositoryMock := entities.NewMockBasketRepository(ctrl)

	// create basket with basket id, otherwise a new basket will be created
	userBasket, err := basketFactory.NewBasketWithID(basketID, userID)
	require.NoError(t, err)

	basketRepositoryMock.EXPECT().FindByUserId(userID).Return(userBasket, nil)

	basket, basketErr := basketFactory.NewBasketWithID(basketID, userID)
	require.NoError(t, basketErr)

	basket.AddItem(product1ID, 1)

	basketRepositoryMock.EXPECT().Save(basket).Return(basketID, nil)

	productRepositoryMock := warehouse.NewMockProductRepository(ctrl)
	// first the stock check in the usecase
	// second in the basket output service
	productRepositoryMock.EXPECT().Find(product1ID).Return(product1, nil).Times(2)

	basketCreatorService := helper.NewBasketCreatorServiceImpl(basketFactory, basketRepositoryMock)
	basketOutputService := helper.NewBasketOutputService(productRepositoryMock)

	useCase := NewAddProductUseCaseImpl(basketCreatorService, basketOutputService, basketRepositoryMock, productRepositoryMock)

	input := &AddProductUseCaseInput{
		UserID:    userID,
		ProductID: product1ID,
		Count:     1,
	}

	// act

	output, err := useCase.Execute(input)

	// assert

	require.NoError(t, err)
	require.NotNil(t, output)
}
