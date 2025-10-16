package usecases

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/entities"
	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/usecases/helper"
	warehouse "github.com/arkadiusjonczek/clean-architecture-go/internal/domain/warehouse/business/entities"
)

func Test_RemoveProductUseCase_NewRemoveProductUseCaseImpl_ReturnsError(t *testing.T) {
	testCases := map[string]struct {
		input *RemoveProductUseCaseInput
	}{
		"input is nil": {
			input: nil,
		},
		"UserID is empty": {
			input: &RemoveProductUseCaseInput{},
		},
		"ProductID is empty": {
			input: &RemoveProductUseCaseInput{
				UserID: "1337",
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

			useCase := NewRemoveProductUseCaseImpl(basketCreatorService, basketOutputService, basketRepositoryMock, productRepositoryMock)

			_, err := useCase.Execute(testCase.input)

			require.Error(t, err)
			require.ErrorContains(t, err, errorString)
		})
	}
}

// TODO: Test also use cases like product not found, product out of stock etc.
func Test_RemoveProductFromBasketUseCase(t *testing.T) {
	// arrange

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	basketID := "12345"
	userID := "1337"

	product1ID := "1"

	basketFactory := entities.NewBasketFactory()

	basketRepositoryMock := entities.NewMockBasketRepository(ctrl)

	// create basket with basket id, otherwise a new basket will be created
	userBasket, err := basketFactory.NewBasketWithID(basketID, userID)
	require.NoError(t, err)

	userBasket.AddItem(product1ID, 1)

	basketRepositoryMock.EXPECT().FindByUserId(userID).Return(userBasket, nil)
	basketRepositoryMock.EXPECT().Save(userBasket).Return(basketID, nil)

	productRepositoryMock := warehouse.NewMockProductRepository(ctrl)

	basketCreatorService := helper.NewBasketCreatorServiceImpl(basketFactory, basketRepositoryMock)

	basketOutputService := helper.NewBasketOutputService(productRepositoryMock)

	useCase := NewRemoveProductUseCaseImpl(basketCreatorService, basketOutputService, basketRepositoryMock, productRepositoryMock)

	input := &RemoveProductUseCaseInput{
		UserID:    userID,
		ProductID: product1ID,
	}

	// act

	output, err := useCase.Execute(input)

	// assert

	require.NoError(t, err)
	require.NotNil(t, output)
}
