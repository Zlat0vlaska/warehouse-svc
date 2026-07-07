package warehouse

import (
	"errors"
	"fmt"
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

		fmt.Printf("Попытка добавить дубликат %s %s\n\n", p.ID, p.Name)

		return errors.New("product already exists\n\n")
	}

	s.products[p.ID] = &p
	fmt.Printf("Добавлен: %s %s\n\n", p.ID, p.Name)

	return nil

}

func (s *Store) Get(id string) (Product, error) {

	if p, ok := s.products[id]; ok {

		fmt.Printf("Получить %s: \n", p.ID)
		fmt.Printf("%+v \n\n", p)
		return *p, nil
	}

	return Product{}, errors.New("product not found\n")
}

func (s *Store) List() []Product {

	var sl = make([]Product, 0, len(s.products))

	fmt.Println("Список продуктов:")
	for _, value := range s.products {
		sl = append(sl, *value)
		fmt.Printf("  - %s %s (осталось: %d)\n", value.ID, value.Name, value.Stock)
	}

	fmt.Printf("\n")
	return sl
}

func (s *Store) UpdateStock(id string, delta int) error {

	p, err := s.Get(id)
	if err != nil {
		fmt.Printf("Обновляем сток %s на %d:\n", p.ID, delta)
		fmt.Printf("ошибка: update stock %s: product not found\n\n", p.ID)

		return errors.New("product not found\n\n")
	}

	p.Stock += delta

	if p.Stock < 0 {
		p.Stock = 0
	}

	fmt.Printf("Обновляем сток %s на %d...\n\n", p.ID, delta)

	return nil
}
