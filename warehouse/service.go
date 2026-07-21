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
	if p.ID == "" || p.Name == "" || p.Price <= 0 || p.Stock < 0 {
		return fmt.Errorf("add product %q: name must not be empty: %w", p.ID, ErrValidation)
	}
	// тут будет валидация: if p.Name == "" { return ... }
	// тут будет обогащение: p.CreatedAt = time.Now()
	// пока — просто прокси, и это нормально
	return svc.repo.Add(p)
}

func (svc *ProductService) Get(id string) (Product, error) {
	return svc.repo.Get(id)
}

func (svc *ProductService) List() []Product {
	return svc.repo.List()
}

func (svc *ProductService) UpdateStock(id string, delta int) error {
	if delta != 0 {
		return svc.repo.UpdateStock(id, delta)
	}
	return fmt.Errorf("update stock for product %q: delta must not be zero: %w", id, ErrValidation)
}
