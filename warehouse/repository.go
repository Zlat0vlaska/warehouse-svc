package warehouse

import (
	"context"
	"fmt"
	"sync"
)

type MemoryRepository struct {
	mu       sync.RWMutex
	products map[string]*Product
}

var _ productRepository = (*MemoryRepository)(nil)

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{

		products: make(map[string]*Product),
	}
}

func (r *MemoryRepository) Add(ctx context.Context, p Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.products[p.ID]; ok {
		return fmt.Errorf("add %q: %w", p.ID, ErrAlreadyExists)
	}
	r.products[p.ID] = &p
	return nil
}

func (r *MemoryRepository) Get(ctx context.Context, id string) (Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if p, ok := r.products[id]; ok {
		return *p, nil
	}
	return Product{}, ErrNotFound
}

func (r *MemoryRepository) List(ctx context.Context) ([]Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var sl = make([]Product, 0, len(r.products))
	for _, value := range r.products {
		sl = append(sl, *value)
	}
	return sl, nil
}

func (r *MemoryRepository) UpdateStock(ctx context.Context, id string, delta int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	p, ok := r.products[id]
	if !ok {
		return ErrNotFound
	}
	if p.Stock+delta < 0 {
		return ErrInsufficientStock
	}
	p.Stock += delta
	return nil
}
