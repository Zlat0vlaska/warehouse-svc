package warehouse

import "errors"

type Product struct {
	ID    string
	Name  string
	Price int
	Stock int
}

var (
	ErrNotFound          = errors.New("product not found")
	ErrAlreadyExists     = errors.New("product already exists")
	ErrInsufficientStock = errors.New("insufficient stock")
)
