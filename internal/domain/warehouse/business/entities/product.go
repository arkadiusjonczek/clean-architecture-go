package entities

type Product struct {
	ID    string
	Name  string
	Price *ProductPrice
	Stock int
}

// ProductPrice is a value object
type ProductPrice struct {
	Value    float64
	Currency string
}
