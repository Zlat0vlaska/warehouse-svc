package warehouse

type ProductService struct {
	repo *MemoryRepository
}

func NewProductService(repo *MemoryRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) AddProduct(p Product) error {
	// тут будет валидация: if p.Name == "" { return ... }
	// тут будет обогащение: p.CreatedAt = time.Now()
	// пока — просто прокси, и это нормально
	return s.repo.Add(p)
}

func (s *ProductService) GetProduct(id string) (Product, error) {

	return s.repo.Get(id)
}

func (s *ProductService) ListProduct() []Product {

	return s.repo.List()
}

func (s *ProductService) UpdateStockProduct(id string, delta int) error {

	return s.repo.UpdateStock(id, delta)
}
