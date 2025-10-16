package usecases

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/entities"
	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/usecases/helper"
	warehouse "github.com/arkadiusjonczek/clean-architecture-go/internal/domain/warehouse/business/entities"
)

func Test_ShowBasketUseCase_NewShowBasketUseCaseImpl_ReturnsError(t *testing.T) {
	testCases := map[string]struct {
		input *ShowBasketUseCaseInput
	}{
		"input is nil": {
			input: nil,
		},
		"UserID is empty": {
			input: &ShowBasketUseCaseInput{},
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

			useCase := NewShowBasketUseCaseImpl(basketCreatorService, basketOutputService)

			_, err := useCase.Execute(testCase.input)

			require.Error(t, err)
			require.ErrorContains(t, err, errorString)
		})
	}
}

func Test_ShowBasketUseCase(t *testing.T) {
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

	userBasket.AddItem(product1ID, 1)

	basketRepositoryMock.EXPECT().FindByUserId(userID).Return(userBasket, nil)

	productRepositoryMock := warehouse.NewMockProductRepository(ctrl)
	productRepositoryMock.EXPECT().Find(product1ID).Return(product1, nil)

	basketCreatorService := helper.NewBasketCreatorServiceImpl(basketFactory, basketRepositoryMock)

	basketOutputService := helper.NewBasketOutputService(productRepositoryMock)

	useCase := NewShowBasketUseCaseImpl(basketCreatorService, basketOutputService)

	input := &ShowBasketUseCaseInput{
		UserID: userID,
	}

	// act

	output, err := useCase.Execute(input)

	// assert

	require.NoError(t, err)
	require.NotNil(t, output)
}
