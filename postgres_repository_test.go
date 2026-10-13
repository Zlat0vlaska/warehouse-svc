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
