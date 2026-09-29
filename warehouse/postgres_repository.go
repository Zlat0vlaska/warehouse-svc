package warehouse

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

var _ productRepository = (*PostgresRepository)(nil)

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Add(ctx context.Context, p Product) error {
	const query = `
        INSERT INTO products (id, name, price, stock)
        VALUES ($1, $2, $3, $4)
    `
	_, err := r.pool.Exec(ctx, query, p.ID, p.Name, p.Price, p.Stock)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("add %q: %w", p.ID, ErrAlreadyExists)
		}
		return fmt.Errorf("add product %q: %w", p.ID, err)
	}
	return err
}

func (r *PostgresRepository) Get(ctx context.Context, id string) (Product, error) {
	const query = `
		SELECT id, name, price, stock
		FROM products
		WHERE id = $1
	`
	var p Product
	err := r.pool.QueryRow(ctx, query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Product{}, fmt.Errorf("get %q: %w", id, ErrNotFound)
		}
		return Product{}, fmt.Errorf("get %q: %w", id, err)
	}
	return p, nil
}

func (r *PostgresRepository) List(ctx context.Context) ([]Product, error) {
	const query = `
		SELECT id, name, price, stock
		FROM products
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}
	defer rows.Close()

	var sl []Product

	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
		if err != nil {
			return nil, fmt.Errorf("scan product: %w", err)
		}
		sl = append(sl, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return sl, nil
}

func (r *PostgresRepository) UpdateStock(ctx context.Context, id string, delta int) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var stock int
	err = tx.QueryRow(ctx,
		"SELECT stock FROM products WHERE id = $1 FOR UPDATE",
		id,
	).Scan(&stock)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("update stock %q: %w", id, ErrNotFound)
		}
		return fmt.Errorf("update  stock %q: %w", id, err)
	}

	if stock+delta < 0 {
		return fmt.Errorf("update stock %q: %w", id, ErrInsufficientStock)
	}

	_, err = tx.Exec(ctx,
		`UPDATE products SET stock = $2 WHERE id = $1`,
		id, stock+delta,
	)

	if err != nil {
		return fmt.Errorf("update stock %q: %w", id, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}
