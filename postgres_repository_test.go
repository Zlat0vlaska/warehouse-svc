//go:build integration
// +build integration

package warehouse

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupPostgres(t *testing.T) *pgxpool.Pool {
	t.Helper()
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("start postgres container: %v", err)
	}
	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Logf("terminate container: %v", err)
		}
	})

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("connection string: %v", err)
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Fatalf("new pool: %v", err)
	}
	t.Cleanup(func() { pool.Close() })

	// прогнать миграцию руками (goose внутри теста — лишнее)
	_, err = pool.Exec(ctx, `
        CREATE TABLE products (
            id         TEXT PRIMARY KEY,
            name       TEXT NOT NULL,
            price      INTEGER NOT NULL,
            stock      INTEGER NOT NULL
        )
    `)
	if err != nil {
		t.Fatalf("create table: %v", err)
	}

	return pool
}

func TestPostgresRepository_Add(t *testing.T) {
	pool := setupPostgres(t)
	repo := NewPostgresRepository(pool)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		err := repo.Add(ctx, Product{ID: "1", Name: "X", Price: 100, Stock: 10})
		if err != nil {
			t.Fatalf("add: %v", err)
		}
	})

	t.Run("duplicate", func(t *testing.T) {
		err := repo.Add(ctx, Product{ID: "1", Name: "Y", Price: 50, Stock: 5})
		if !errors.Is(err, ErrAlreadyExists) {
			t.Fatalf("got %v, want ErrAlreadyExists", err)
		}
	})
}

func TestPostgresRepository_Get(t *testing.T) {
	pool := setupPostgres(t)
	repo := NewPostgresRepository(pool)
	ctx := context.Background()

	existingProduct := Product{ID: "1", Name: "X", Price: 100, Stock: 10}
	err := repo.Add(ctx, existingProduct)
	if err != nil {
		t.Fatalf("failed to seed data: %v", err)
	}

	t.Run("success", func(t *testing.T) {
		p, err := repo.Get(ctx, "1")
		if err != nil {
			t.Fatalf("get: %v", err)
		}

		// Дополнительно проверяем, что база вернула именно то, что мы сохранили
		if p != existingProduct {
			t.Errorf("got %+v, want %+v", p, existingProduct)
		}
	})

	t.Run("notFound", func(t *testing.T) {
		_, err := repo.Get(ctx, "100")
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("got %v, want ErrNotFound", err)
		}
	})
}

func TestPostgresRepository_List(t *testing.T) {
	pool := setupPostgres(t)
	repo := NewPostgresRepository(pool)
	ctx := context.Background()

	products := []Product{
		{ID: "1", Name: "Product A", Price: 100, Stock: 10},
		{ID: "2", Name: "Product B", Price: 200, Stock: 5},
		{ID: "3", Name: "Product C", Price: 300, Stock: 0},
	}

	for _, product := range products {
		if err := repo.Add(ctx, product); err != nil {
			t.Fatalf("failed to seed product %d: %v", product.ID, err)
		}

	}

	t.Run("success", func(t *testing.T) {
		gotProducts, err := repo.List(ctx)
		if err != nil {
			t.Fatalf("list: %v", err)
		}

		if len(gotProducts) != len(products) {
			t.Fatalf("got %v, want %v", gotProducts, products)
		}

		for i := range gotProducts {
			if gotProducts[i] != products[i] {
				t.Fatalf("got %v, want %v", gotProducts[i], products[i])
			}
		}
	})
}

func TestPostgresRepository_UpdateStock(t *testing.T) {
	pool := setupPostgres(t)
	repo := NewPostgresRepository(pool)
	ctx := context.Background()

	existingProduct := Product{ID: "1", Name: "X", Price: 100, Stock: 10}
	err := repo.Add(ctx, existingProduct)
	if err != nil {
		t.Fatalf("failed to seed data: %v", err)
	}

	t.Run("success", func(t *testing.T) {
		err := repo.UpdateStock(ctx, existingProduct.id, existingProduct.delta)
		if err != nil {
			t.Fatalf("failed to seed data: %v", err)
		}
		existingProduct.Stock += existingProduct.delta
		p, err := repo.Get(ctx, existingProduct.id)
		if err != nil {
			t.Fatalf("get: %v", err)
		}
		if p.Stock != existingProduct.Stock {
			t.Fatalf("got stock %d, want %d", p.Stock, existingProduct.Stock)
		}
	})
	t.Run("NotFound", func(t *testing.T) {
		err := repo.UpdateStock(ctx, "nonexistent", 5)
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("got %v, want ErrNotFound", err)
		}
	})

	t.Run("NotEnoughStock", func(t *testing.T) {
		err := repo.UpdateStock(ctx, existingProduct.id, -100)
		if !errors.Is(err, ErrInsufficientStock) {
			t.Fatalf("got %v, want ErrInsufficientStock", err)
		}
	})
}
