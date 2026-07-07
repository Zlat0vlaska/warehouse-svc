package main

import "github.com/Zlat0vlaska/warehouse-svc/warehouse"

func main() {
	store := warehouse.New()

	store.Add(warehouse.Product{ID: "1", Name: "Product 1", Price: 100, Stock: 10})
	store.Add(warehouse.Product{ID: "2", Name: "Product 2", Price: 200, Stock: 20})
	store.Add(warehouse.Product{ID: "3", Name: "Product 3", Price: 300, Stock: 0})

	store.List()

	store.Get("1")
	store.Get("4")

	store.Add(warehouse.Product{ID: "1", Name: "Product 1", Price: 100, Stock: 10})

	store.UpdateStock("2", 15)
	store.UpdateStock("4", 15)

	store.List()
}
