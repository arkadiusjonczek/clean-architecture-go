package helper

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_ProductPriceSimulatorBackgroundServiceImpl_NewProductPriceSimulator_ReturnsError(t *testing.T) {
	service, err := NewProductPriceSimulatorBackgroundService(nil)

	require.Error(t, err)
	require.Nil(t, service)
}

func Test_ProductPriceSimulatorBackgroundServiceImpl_NewProductPriceSimulator(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockProductPriceSimulatorService(ctrl)

	service, err := NewProductPriceSimulatorBackgroundService(mock)

	require.NoError(t, err)
	require.NotNil(t, service)
}
