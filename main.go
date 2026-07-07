package main

import (
	"fmt"

	"github.com/Zlat0vlaska/warehouse-svc/warehouse"
)

func main() {
	store := warehouse.New()

	if err := store.Add(warehouse.Product{ID: "1", Name: "Product 1", Price: 100, Stock: 10}); err != nil {
		fmt.Printf("Ошибка добавления: %v\n", err)
	} else {
		fmt.Println("Продукт 1 успешно добавлен")
	}

	if err := store.Add(warehouse.Product{ID: "2", Name: "Product 2", Price: 200, Stock: 20}); err != nil {
		fmt.Printf("Ошибка добавления: %v\n", err)
	} else {
		fmt.Println("Продукт 2 успешно добавлен")
	}

	store.List()
	if p, err := store.Get("1"); err != nil {
		fmt.Printf("Ошибка получения: %v\n", err)
	} else {
		fmt.Printf("Получен продукт: %+v\n", p)
	}

	if p, err := store.Get("4"); err != nil {
		fmt.Printf("Ошибка получения: %v\n", err)
	} else {
		fmt.Printf("Получен продукт: %+v\n", p)
	}

	if err := store.Add(warehouse.Product{ID: "1", Name: "Product 1", Price: 100, Stock: 10}); err != nil {
		fmt.Printf("Ошибка добавления дубликата: %v\n", err)
	} else {
		fmt.Println("Дубликат добавлен (не должно случиться)")
	}

	if err := store.UpdateStock("2", 15); err != nil {
		fmt.Printf("Ошибка обновления стока: %v\n", err)
	} else {
		fmt.Println("Сток продукта 2 обновлён")
	}

	// Проверим обновлённый сток
	if p, err := store.Get("2"); err != nil {
		fmt.Printf("Ошибка получения: %v\n", err)
	} else {
		fmt.Printf("Продукт 2 после обновления: %+v\n", p)
	}

	// Обновление несуществующего
	if err := store.UpdateStock("4", 15); err != nil {
		fmt.Printf("Ошибка обновления стока: %v\n", err)
	} else {
		fmt.Println("Сток продукта 4 обновлён")
	}

	// --- Снова список ---
	store.List()
}
