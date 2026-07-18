package main

import (
	"log"
	"net/http"

	"github.com/Zlat0vlaska/warehouse-svc/warehouse"
)

func main() {
	repo := warehouse.NewMemoryRepository()
	svc := warehouse.NewProductService(repo)
	_ = svc.Add(warehouse.Product{ID: "1", Name: "Coffee", Price: 500, Stock: 10})
	_ = svc.Add(warehouse.Product{ID: "2", Name: "Tea", Price: 300, Stock: 20})
	mux := http.NewServeMux()
	warehouse.RegisterRoutes(mux, svc)
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
