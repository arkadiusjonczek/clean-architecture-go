package inmemory

import (
	"fmt"

	"github.com/arkadiusjonczek/clean-architecture-go/internal/domain/warehouse"
)

var _ warehouse.ProductRepository = (*InMemoryProductRepository)(nil)

type InMemoryProductRepository struct {
	products map[string]*warehouse.Product
}

func NewInMemoryProductRepository() warehouse.ProductRepository {
	return &InMemoryProductRepository{
		products: make(map[string]*warehouse.Product),
	}
}

func (repository *InMemoryProductRepository) Save(product *warehouse.Product) {
	repository.products[product.ID] = product
}

func (repository *InMemoryProductRepository) Find(id string) (*warehouse.Product, error) {
	product, productExists := repository.products[id]
	if !productExists {
		return nil, fmt.Errorf("product not found")
	}

	return product, nil
}
