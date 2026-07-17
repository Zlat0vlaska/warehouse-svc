package main

import (
	"encoding/json"
	"errors"
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

type updateStockRequest struct {
	Delta int `json:"delta"`
}

func createProductHandler(store *warehouse.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createProductRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json body", http.StatusBadRequest)
			return
		}

		p := warehouse.Product{
			ID:    req.ID,
			Name:  req.Name,
			Price: req.Price,
			Stock: req.Stock,
		}

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
	}
}

func getProductHandler(store *warehouse.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}

		product, err := store.Get(id)
		if err != nil {
			if errors.Is(err, warehouse.ErrNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(product)
	}
}

func listProductsHandler(store *warehouse.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		list := store.List()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(list)
	}
}

func updateStockHandler(store *warehouse.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}

		var req updateStockRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json body", http.StatusBadRequest)
			return
		}
		err := store.UpdateStock(id, req.Delta)
		if err != nil {
			if errors.Is(err, warehouse.ErrInsufficientStock) {
				http.Error(w, err.Error(), http.StatusConflict)
				return
			}
		}
		product, err := store.Get(id)
		if err != nil {
			if errors.Is(err, warehouse.ErrNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(product)
	}
}

func main() {

	store := warehouse.New()

	_ = store.Add(warehouse.Product{ID: "1", Name: "Coffee", Price: 500, Stock: 10})
	_ = store.Add(warehouse.Product{ID: "2", Name: "Tea", Price: 300, Stock: 20})

	mux := http.NewServeMux()
	mux.HandleFunc("GET /products", listProductsHandler(store))
	mux.HandleFunc("GET /products/{id}", getProductHandler(store))
	mux.HandleFunc("POST /products", createProductHandler(store))
	mux.HandleFunc("PATCH /products/{id}/stock", updateStockHandler(store))

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
