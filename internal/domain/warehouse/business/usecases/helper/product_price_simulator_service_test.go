package helper

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/warehouse/business/entities"
)

func Test_ProductPriceSimulatorServiceImpl_NewProductPriceSimulator(t *testing.T) {
	service, err := NewProductPriceSimulatorService(nil)

	require.Error(t, err)
	require.Nil(t, service)
}

func Test_ProductPriceSimulatorServiceImpl_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	products := []*entities.Product{
		{
			ID:    "A12345",
			Name:  "Product A12345",
			Stock: 10,
			Price: &entities.ProductPrice{
				Value:    13.37,
				Currency: "EUR",
			},
		},
	}

	mockProductRepository := entities.NewMockProductRepository(ctrl)
	mockProductRepository.EXPECT().FindAll().Return(products).Times(1)
	mockProductRepository.EXPECT().Save(products[0]).Return().Times(1)

	service, err := NewProductPriceSimulatorService(mockProductRepository)

	require.NoError(t, err)
	require.NotNil(t, service)

	oldPrice := products[0].Price.Value

	service.Execute()

	require.NotEqual(t, oldPrice, products[0].Price.Value)
}
