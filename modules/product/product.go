package product

// Product represents an item that can be offered or sold.
type Product struct {
	ID          int64
	Name        string
	Description string
	PriceCents  int64
}

// Service defines product-related business operations.
type Service interface {
	Create(p Product) (Product, error)
	GetByID(id int64) (Product, error)
	List() ([]Product, error)
}
