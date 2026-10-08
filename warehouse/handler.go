package warehouse

import (
	"encoding/json"
	"errors"
	"log"
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

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("encode response: %v", err)
	}
}

func createProductHandler(svc *ProductService) http.HandlerFunc {
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

		if err := svc.Add(r.Context(), p); err != nil {
			switch {
			case errors.Is(err, ErrValidation):
				http.Error(w, err.Error(), http.StatusBadRequest)
			case errors.Is(err, ErrAlreadyExists):
				http.Error(w, err.Error(), http.StatusConflict)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		writeJSON(w, http.StatusCreated, p)
	}
}

func getProductHandler(svc *ProductService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}

		product, err := svc.Get(r.Context(), id)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, product)
	}
}

func listProductsHandler(svc *ProductService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		list, err := svc.List(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		writeJSON(w, http.StatusOK, list)
	}
}

func updateStockHandler(svc *ProductService) http.HandlerFunc {
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

		if err := svc.UpdateStock(r.Context(), id, req.Delta); err != nil {
			switch {
			case errors.Is(err, ErrValidation):
				http.Error(w, err.Error(), http.StatusBadRequest)
			case errors.Is(err, ErrNotFound):
				http.Error(w, err.Error(), http.StatusNotFound)
			case errors.Is(err, ErrInsufficientStock):
				http.Error(w, err.Error(), http.StatusConflict)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		product, err := svc.Get(r.Context(), id)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		writeJSON(w, http.StatusOK, product)
	}
}
