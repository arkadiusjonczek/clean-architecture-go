package usecases

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/entities"
	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/usecases/helper"
	warehouse "github.com/arkadiusjonczek/clean-architecture-go/internal/domain/warehouse/business/entities"
)

func Test_UpdateProductCount_NewUpdateProductCountImpl_ReturnsError(t *testing.T) {
	testCases := map[string]struct {
		input *UpdateProductCountUseCaseInput
	}{
		"input is nil": {
			input: nil,
		},
		"input parameter UserID is empty": {
			input: &UpdateProductCountUseCaseInput{},
		},
		"input parameter ProductID is empty": {
			input: &UpdateProductCountUseCaseInput{
				UserID: "1337",
			},
		},
		"input parameter Count is invalid (must be greater than 0)": {
			input: &UpdateProductCountUseCaseInput{
				UserID:    "1337",
				ProductID: "1337",
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

			useCase := NewUpdateProductCountImpl(basketCreatorService, basketOutputService, basketRepositoryMock, productRepositoryMock)

			_, err := useCase.Execute(testCase.input)

			require.Error(t, err)
			require.ErrorContains(t, err, errorString)
		})
	}
}
