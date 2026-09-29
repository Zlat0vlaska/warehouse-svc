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

func (svc *MemoryRepository) Add(ctx context.Context, p Product) error {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	if _, ok := svc.products[p.ID]; ok {
		return fmt.Errorf("add %q: %w", p.ID, ErrAlreadyExists)
	}
	svc.products[p.ID] = &p
	return nil
}

func (svc *MemoryRepository) Get(ctx context.Context, id string) (Product, error) {
	svc.mu.RLock()
	defer svc.mu.RUnlock()

	if p, ok := svc.products[id]; ok {
		return *p, nil
	}
	return Product{}, ErrNotFound
}

func (svc *MemoryRepository) List(ctx context.Context) ([]Product, error) {
	var sl = make([]Product, 0, len(svc.products))
	for _, value := range svc.products {
		sl = append(sl, *value)
	}
	return sl, nil
}

func (svc *MemoryRepository) UpdateStock(ctx context.Context, id string, delta int) error {
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
