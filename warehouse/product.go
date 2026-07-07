package warehouse

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound          = errors.New("...")
	ErrAlreadyExists     = errors.New("...")
	ErrInsufficientStock = errors.New("...")
)

type Product struct {
	ID string

	Name string

	Price int

	Stock int
}

type Store struct {
	products map[string]*Product
}

func New() *Store {

	return &Store{

		products: make(map[string]*Product),
	}

}

func (s *Store) Add(p Product) error {

	if _, ok := s.products[p.ID]; ok {

		return fmt.Errorf("add %q: %w", p.ID, ErrAlreadyExists)
	}

	s.products[p.ID] = &p

	return nil

}

func (s *Store) Get(id string) (Product, error) {

	if p, ok := s.products[id]; ok {

		return *p, nil
	}

	return Product{}, ErrNotFound
}

func (s *Store) List() []Product {

	var sl = make([]Product, 0, len(s.products))

	for _, value := range s.products {
		sl = append(sl, *value)
	}

	return sl
}

func (s *Store) UpdateStock(id string, delta int) error {

	p, ok := s.products[id]

	if !ok {
		return ErrNotFound
	}

	if p.Stock < 0 {
		return ErrInsufficientStock
	}

	p.Stock += delta

	return nil
}
