package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Zlat0vlaska/warehouse-svc/warehouse"
	"github.com/jackc/pgx/v5/pgxpool"
)

func run() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://warehouse:warehouse@localhost:5432/warehouse?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return fmt.Errorf("connect to postgres: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		return fmt.Errorf("ping postgres: %v", err)
	}
	repo := warehouse.NewPostgresRepository(pool)
	svc := warehouse.NewProductService(repo)

	mux := http.NewServeMux()
	warehouse.RegisterRoutes(mux, svc)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,  // сколько ждать заголовков
		ReadTimeout:       10 * time.Second, // всё тело запроса должно прийти за это время
		WriteTimeout:      10 * time.Second, // ответ должен уйти за это время
		IdleTimeout:       60 * time.Second, // keep-alive соединения без активности
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	serverErr := make(chan error, 1)
	go func() {
		log.Println("listening on :8080")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
		close(serverErr)
	}()

	select {
	case err := <-serverErr:
		return fmt.Errorf("server error: %v", err)
	case sig := <-shutdown:
		log.Printf("shutdown signal received: %v", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server shutdown: %v", err)
	}

	log.Println("shutdown complete")

	return nil
}

func main() {

	if err := run(); err != nil {
		log.Fatalf("run: %v", err)
	}

}
