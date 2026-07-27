package warehouse

import (
	"context"
	"fmt"
)

type productRepository interface {
	Add(ctx context.Context, p Product) error
	Get(ctx context.Context, id string) (Product, error)
	List(ctx context.Context) ([]Product, error)
	UpdateStock(ctx context.Context, id string, delta int) error
}
type ProductService struct {
	repo productRepository
}

func NewProductService(repo productRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (svc *ProductService) Add(ctx context.Context, p Product) error {
	if p.ID == "" {
		return fmt.Errorf("add product: id must not be empty: %w", ErrValidation)
	}
	if p.Name == "" {
		return fmt.Errorf("add product %q: name must not be empty: %w", p.ID, ErrValidation)
	}
	if p.Price <= 0 {
		return fmt.Errorf("add product %q: price must be positive, got %d: %w", p.ID, p.Price, ErrValidation)
	}
	if p.Stock < 0 {
		return fmt.Errorf("add product %q: stock must not be negative, got %d: %w", p.ID, p.Stock, ErrValidation)
	}
	return svc.repo.Add(ctx, p)
}

func (svc *ProductService) Get(ctx context.Context, id string) (Product, error) {
	return svc.repo.Get(ctx, id)
}

func (svc *ProductService) List(ctx context.Context) ([]Product, error) {
	return svc.repo.List(ctx)
}

func (svc *ProductService) UpdateStock(ctx context.Context, id string, delta int) error {
	if delta == 0 {
		return fmt.Errorf("update stock for product %q: delta must not be zero: %w", id, ErrValidation)
	}
	return svc.repo.UpdateStock(ctx, id, delta)
}
