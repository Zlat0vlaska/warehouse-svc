package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/Zlat0vlaska/warehouse-svc/warehouse"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://warehouse:warehouse@localhost:5432/warehouse?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("ping postgres: %v", err)
	}
	defer pool.Close()

	repo := warehouse.NewPostgresRepository(pool)
	svc := warehouse.NewProductService(repo)

	mux := http.NewServeMux()
	warehouse.RegisterRoutes(mux, svc)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
