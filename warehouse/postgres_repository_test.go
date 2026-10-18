//go:build integration
// +build integration

package warehouse

import (
	"context"
	"errors"
	"sync"
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
            ID         TEXT PRIMARY KEY,
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
			t.Fatalf("failed to seed product %s: %v", product.ID, err)
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

		want := make(map[string]Product, len(products))
		for _, p := range products {
			want[p.ID] = p
		}

		for _, got := range gotProducts {
			expected, ok := want[got.ID]
			if !ok {
				t.Errorf("unexpected product ID %s", got.ID)
				continue
			}
			if got != expected {
				t.Fatalf("product %s: got %+v, want %+v", got.ID, got, expected)
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
		err := repo.UpdateStock(ctx, existingProduct.ID, 5)
		if err != nil {
			t.Fatalf("failed to seed data: %v", err)
		}
		existingProduct.Stock += 5
		p, err := repo.Get(ctx, existingProduct.ID)
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
		err := repo.UpdateStock(ctx, existingProduct.ID, -100)
		if !errors.Is(err, ErrInsufficientStock) {
			t.Fatalf("got %v, want ErrInsufficientStock", err)
		}
	})
}

func TestPostgresRepository_UpdateStock_Concurrent(t *testing.T) {
	pool := setupPostgres(t)
	repo := NewPostgresRepository(pool)
	ctx := context.Background()

	if err := repo.Add(ctx, Product{ID: "1", Name: "X", Price: 100, Stock: 100}); err != nil {
		t.Fatalf("seed: %v", err)
	}

	const workers = 50
	var wg sync.WaitGroup
	errCh := make(chan error, workers)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := repo.UpdateStock(ctx, "1", -1); err != nil {
				errCh <- err
			}
		}()
	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		t.Errorf("concurrent UpdateStock: %v", err)
	}

	p, _ := repo.Get(ctx, "1")
	wantStock := 100 - workers
	if p.Stock != wantStock {
		t.Fatalf("stock after concurrent updates: got %d, want %d", p.Stock, wantStock)
	}
}
