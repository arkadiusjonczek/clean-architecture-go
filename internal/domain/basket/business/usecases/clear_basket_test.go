package usecases

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/entities"
	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/basket/business/usecases/helper"
	warehouse "github.com/arkadiusjonczek/clean-architecture-go/internal/domain/warehouse/business/entities"
)

func Test_ClearBasketUseCase_NewClearBasketUseCaseImpl_ReturnsError(t *testing.T) {
	testCases := map[string]struct {
		input *ClearBasketUseCaseInput
	}{
		"input is nil": {
			input: nil,
		},
		"UserID is empty": {
			input: &ClearBasketUseCaseInput{},
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

			useCase := NewClearBasketUseCaseImpl(basketCreatorService, basketOutputService, basketRepositoryMock)

			_, err := useCase.Execute(testCase.input)

			require.Error(t, err)
			require.ErrorContains(t, err, errorString)
		})
	}
}
