package entities

//go:generate mockgen -source=basket_repository.go -destination=basket_repository_mock.go -package=entities

type BasketRepository interface {
	Find(id string) (*Basket, error)
	FindByUserId(userId string) (*Basket, error) // special function
	Save(basket *Basket) (string, error)
}

var _ error = (*BasketNotFoundError)(nil)

type BasketNotFoundError struct {
}

func (err *BasketNotFoundError) Error() string {
	return "basket not found"
}
