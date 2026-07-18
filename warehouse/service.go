package warehouse

type ProductService struct {
	repo *MemoryRepository
}

func NewProductService(repo *MemoryRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (svc *ProductService) Add(p Product) error {
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
	return svc.repo.UpdateStock(id, delta)
}
