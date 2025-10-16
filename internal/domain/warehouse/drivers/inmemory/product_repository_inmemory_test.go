package inmemory

import (
	"testing"

	"github.com/stretchr/testify/require"

	warehouse "github.com/arkadiusjonczek/clean-architecture-go/internal/domain/warehouse/business/entities"
)

func Test_InMemoryProductRepository_Find(t *testing.T) {
	repository := NewInMemoryProductRepository()

	require.NotNil(t, repository)

	productID := "A12345"
	productName := "Product A12345"

	product, err := repository.Find(productID)

	require.Error(t, err)
	require.Nil(t, product)

	repository.Save(&warehouse.Product{
		ID:   productID,
		Name: productName,
	})

	product, err = repository.Find(productID)

	require.NoError(t, err)
	require.NotNil(t, product)
	require.Equal(t, productID, product.ID)
	require.Equal(t, productName, product.Name)
}
func Test_InMemoryProductRepository_FindAll(t *testing.T) {
	repository := NewInMemoryProductRepository()

	require.NotNil(t, repository)

	require.NotNil(t, repository.FindAll())
	require.Empty(t, repository.FindAll())
}
