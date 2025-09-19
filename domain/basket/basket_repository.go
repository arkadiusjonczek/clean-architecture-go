package basket

//go:generate mockgen -source=basket_repository.go -destination=basket_repository_mock.go -package=basket

type BasketRepository interface {
	Find(id string) (*Basket, error)
	FindByUserId(userId string) (*Basket, error)
	Save(basket *Basket) (string, error)
}
