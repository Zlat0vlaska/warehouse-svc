package warehouse

import (
	"context"
	"log"
	"time"
)

type LoggingRepository struct {
	repo   productRepository
	logger *log.Logger
}

func NewLoggingRepository(repo productRepository, logger *log.Logger) *LoggingRepository {
	return &LoggingRepository{
		repo:   repo,
		logger: logger,
	}
}

func (l *LoggingRepository) Add(ctx context.Context, p Product) error {
	start := time.Now()
	err := l.repo.Add(ctx, p)
	l.logger.Printf("Add(%+v) -> err: %v, took: %v", p, err, time.Since(start))
	return err
}

func (l *LoggingRepository) Get(ctx context.Context, id string) (Product, error) {
	start := time.Now()
	p, err := l.repo.Get(ctx, id)
	if err != nil {
		l.logger.Printf("Get(%s) -> err: %v, took: %v", id, err, time.Since(start))
	} else {
		l.logger.Printf("Get(%s) -> found: %+v, took: %v", id, p, time.Since(start))
	}
	return p, err
}

func (l *LoggingRepository) List(ctx context.Context) ([]Product, error) {
	start := time.Now()
	products, err := l.repo.List(ctx)
	if err != nil {
		l.logger.Printf("List() -> err: %v, took: %v", err, time.Since(start))
	} else {
		l.logger.Printf("List() -> count: %d, took: %v", len(products), time.Since(start))
	}
	return products, err
}

func (l *LoggingRepository) UpdateStock(ctx context.Context, id string, delta int) error {
	start := time.Now()
	err := l.repo.UpdateStock(ctx, id, delta)
	l.logger.Printf("UpdateStock(%s, %d) -> err: %v, took: %v", id, delta, err, time.Since(start))
	return err
}
