package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Zlat0vlaska/warehouse-svc/warehouse"
)

type createProductRequest struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

func createProductHandler(store *warehouse.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createProductRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		p := warehouse.Product{
			ID:    req.ID,
			Name:  req.Name,
			Price: req.Price,
			Stock: req.Stock,
		}

		switch r.Method {
		case http.MethodGet:
			if _, err := store.Get(p.ID); err != nil {
				if errors.Is(err, warehouse.ErrNotFound) {
					http.Error(w, err.Error(), http.StatusNotFound)
					return
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusAccepted)
			json.NewEncoder(w).Encode(p)

		case http.MethodPost:
			if err := store.Add(p); err != nil {
				if errors.Is(err, warehouse.ErrAlreadyExists) {
					http.Error(w, err.Error(), http.StatusConflict)
					return
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(p)
		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}

	}
}

func main() {
	mux := http.NewServeMux()
	store := warehouse.New()
	log.Fatal(http.ListenAndServe(":8080", mux))

	mux.HandleFunc("GET /products", listHandler)
	mux.HandleFunc("GET /products/{id}", getHandler)
	mux.HandleFunc("POST /products", createHandler)
	mux.HandleFunc("PATCH /products/{id}/stock", updateStockHandler)

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
	if err := store.Add(warehouse.Product{ID: "3", Name: "Product 3", Price: 300, Stock: 30}); err != nil {
		fmt.Printf("Ошибка добавления: %v\n", err)
	} else {
		fmt.Println("Продукт 3 успешно добавлен")
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
	if p, err := store.Get("2"); err != nil {
		fmt.Printf("Ошибка получения: %v\n", err)
	} else {
		fmt.Printf("Продукт 2 после обновления: %+v\n", p)
	}
	if err := store.UpdateStock("4", 15); err != nil {
		fmt.Printf("Ошибка обновления стока: %v\n", err)
	} else {
		fmt.Println("Сток продукта 4 обновлён")
	}
	store.List()
	if err := store.UpdateStock("3", -100); err != nil {
		fmt.Printf("Ошибка обновления стока: %v\n", err)
	} else {
		fmt.Println("Сток продукта 3 обновлён")
	}
	store.List()
}
