package dto

type BasketDTO struct {
	Items []*BasketItem
}

type BasketItem struct {
	Product *Product
	Count   int
}

type Product struct {
	ID    string
	Name  string
	Price *ProductPrice
}

type ProductPrice struct {
	Value    string
	Currency string
}
