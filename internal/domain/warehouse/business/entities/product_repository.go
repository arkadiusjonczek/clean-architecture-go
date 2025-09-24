package entities

//go:generate mockgen -source=product_repository.go -destination=product_repository_mock.go -package=entities

type ProductRepository interface {
	Find(id string) (*Product, error)
	FindAll() []*Product
	Save(product *Product)
}
