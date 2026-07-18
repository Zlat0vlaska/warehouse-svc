package warehouse

import "fmt"

type MemoryRepository struct {
	products map[string]*Product
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		products: make(map[string]*Product),
	}
}

func (s *MemoryRepository) Add(p Product) error {
	if _, ok := s.products[p.ID]; ok {
		return fmt.Errorf("add %q: %w", p.ID, ErrAlreadyExists)
	}
	s.products[p.ID] = &p
	return nil
}

func (s *MemoryRepository) Get(id string) (Product, error) {
	if p, ok := s.products[id]; ok {
		return *p, nil
	}
	return Product{}, ErrNotFound
}

func (s *MemoryRepository) List() []Product {
	var sl = make([]Product, 0, len(s.products))
	for _, value := range s.products {
		sl = append(sl, *value)
	}
	return sl
}

func (s *MemoryRepository) UpdateStock(id string, delta int) error {
	p, ok := s.products[id]
	if !ok {
		return ErrNotFound
	}
	if p.Stock+delta < 0 {
		return ErrInsufficientStock
	}
	p.Stock += delta
	return nil
}
