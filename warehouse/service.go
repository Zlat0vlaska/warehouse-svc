package warehouse

import "fmt"

type productRepository interface {
	Add(p Product) error
	Get(id string) (Product, error)
	List() []Product
	UpdateStock(id string, delta int) error
}
type ProductService struct {
	repo productRepository
}

func NewProductService(repo productRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (svc *ProductService) Add(p Product) error {
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
	return svc.repo.Add(p)
	// тут будет валидация: if p.Name == "" { return ... }
	// тут будет обогащение: p.CreatedAt = time.Now()
	// пока — просто прокси, и это нормально
}

func (svc *ProductService) Get(id string) (Product, error) {
	return svc.repo.Get(id)
}

func (svc *ProductService) List() []Product {
	return svc.repo.List()
}

func (svc *ProductService) UpdateStock(id string, delta int) error {
	if delta == 0 {
		return fmt.Errorf("update stock for product %q: delta must not be zero: %w", id, ErrValidation)
	}
	return svc.repo.UpdateStock(id, delta)

}
