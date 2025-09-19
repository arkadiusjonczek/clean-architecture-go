package warehouse

//go:generate mockgen -source=product_repository.go -destination=product_repository_mock.go -package=warehouse

type ProductRepository interface {
	Save(product *Product)
	Find(id string) (*Product, error)
}
