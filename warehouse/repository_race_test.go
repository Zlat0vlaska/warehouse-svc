package warehouse

import (
	"context"
	"fmt"
	"sync"
	"testing"
)

func TestMemoryRepository_ConcurrentWrites(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_ = repo.Add(ctx, Product{
				ID:    string(rune('A'+i%26)) + string(rune('0'+i/26)),
				Name:  "X",
				Price: 1,
				Stock: 1,
			})
		}(i)

		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, _ = repo.Get(ctx, fmt.Sprintf("seed-%d", i%10))
		}(i)

		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = repo.List(ctx)
		}()

		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_ = repo.UpdateStock(ctx, fmt.Sprintf("seed-%d", i%10), -1)
		}(i)

	}
	wg.Wait()
}
