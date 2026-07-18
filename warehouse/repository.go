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

func (svc *MemoryRepository) Add(p Product) error {
	if _, ok := svc.products[p.ID]; ok {
		return fmt.Errorf("add %q: %w", p.ID, ErrAlreadyExists)
	}
	svc.products[p.ID] = &p
	return nil
}

func (svc *MemoryRepository) Get(id string) (Product, error) {
	if p, ok := svc.products[id]; ok {
		return *p, nil
	}
	return Product{}, ErrNotFound
}

func (svc *MemoryRepository) List() []Product {
	var sl = make([]Product, 0, len(svc.products))
	for _, value := range svc.products {
		sl = append(sl, *value)
	}
	return sl
}

func (svc *MemoryRepository) UpdateStock(id string, delta int) error {
	p, ok := svc.products[id]
	if !ok {
		return ErrNotFound
	}
	if p.Stock+delta < 0 {
		return ErrInsufficientStock
	}
	p.Stock += delta
	return nil
}
