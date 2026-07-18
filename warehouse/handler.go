package warehouse

import (
	"encoding/json"
	"errors"
	"net/http"
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

func RegisterRoutes(mux *http.ServeMux, svc *ProductService) {

	mux.HandleFunc("GET /products", listProductsHandler(svc))
	mux.HandleFunc("GET /products/{id}", getProductHandler(svc))
	mux.HandleFunc("POST /products", createProductHandler(svc))
	mux.HandleFunc("PATCH /products/{id}/stock", updateStockHandler(svc))
}

func createProductHandler(store *ProductService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createProductRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json body", http.StatusBadRequest)
			return
		}

		p := Product{
			ID:    req.ID,
			Name:  req.Name,
			Price: req.Price,
			Stock: req.Stock,
		}

		if err := store.AddProduct(p); err != nil {
			if errors.Is(err, ErrAlreadyExists) {
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

func getProductHandler(s *ProductService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}

		product, err := s.GetProduct(id)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
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

func listProductsHandler(s *ProductService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		list := s.ListProduct()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(list)
	}
}

func updateStockHandler(s *ProductService) http.HandlerFunc {
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

		if err := s.UpdateStockProduct(id, req.Delta); err != nil {
			switch {
			case errors.Is(err, ErrNotFound):
				http.Error(w, err.Error(), http.StatusNotFound)
			case errors.Is(err, ErrInsufficientStock):
				http.Error(w, err.Error(), http.StatusConflict)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		product, err := s.GetProduct(id)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
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
