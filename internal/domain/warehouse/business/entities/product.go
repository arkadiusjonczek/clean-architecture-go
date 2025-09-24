package entities

// ProductPrice is a value object
type ProductPrice struct {
	Price    float64
	Currency string
}

type Product struct {
	ID    string
	Name  string
	Price *ProductPrice
	Stock int
}
