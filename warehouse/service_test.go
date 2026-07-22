package warehouse

import (
	"errors"
	"strings"
	"testing"
)

type mockRepo struct {
	addFn         func(p Product) error
	getFn         func(id string) (Product, error)
	listFn        func() []Product
	updateStockFn func(id string, delta int) error
}

func (m *mockRepo) Add(p Product) error                    { return m.addFn(p) }
func (m *mockRepo) Get(id string) (Product, error)         { return m.getFn(id) }
func (m *mockRepo) List() []Product                        { return m.listFn() }
func (m *mockRepo) UpdateStock(id string, delta int) error { return m.updateStockFn(id, delta) }

func TestProductServiceAdd(t *testing.T) {
	tests := []struct {
		name          string
		input         Product
		repoErr       error  // что вернёт mockRepo.Add
		wantErr       error  // какая ошибка ожидается (nil = успех)
		wantMsgSubstr string // если не пусто, то проверяем, что сообщение об ошибке содержит эту подстроку
	}{
		{
			name:    "success",
			input:   Product{ID: "1", Name: "X", Price: 100, Stock: 10},
			repoErr: nil,
			wantErr: nil,
		},
		{
			name:          "empty name",
			input:         Product{ID: "1", Name: "", Price: 100, Stock: 10},
			wantErr:       ErrValidation, // до репо не должны дойти
			wantMsgSubstr: "name must not be empty",
		},
		{
			name:          "negative price",
			input:         Product{ID: "1", Name: "X", Price: -1, Stock: 10},
			wantErr:       ErrValidation,
			wantMsgSubstr: "price must be positive",
		},
		{
			name:          "negative stock",
			input:         Product{ID: "1", Name: "X", Price: 100, Stock: -1},
			wantErr:       ErrValidation,
			wantMsgSubstr: "stock must not be negative",
		},
		{
			name:          "negative stock error",
			input:         Product{ID: "1", Name: "X", Price: 100, Stock: -1},
			wantErr:       ErrValidation,
			wantMsgSubstr: "name must not be empty",
		},
		{

			name:    "repo says duplicate",
			input:   Product{ID: "1", Name: "X", Price: 100, Stock: 10},
			repoErr: ErrAlreadyExists,
			wantErr: ErrAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepo{
				addFn: func(p Product) error { return tt.repoErr },
			}
			svc := NewProductService(repo)

			err := svc.Add(tt.input)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("got err %v, want %v", err, tt.wantErr)
			}
			if tt.wantMsgSubstr != "" && !strings.Contains(err.Error(), tt.wantMsgSubstr) {
				t.Errorf("err message %q does not contain %q", err.Error(), tt.wantMsgSubstr)
			}
		})
	}
}

func TestProductServiceUpdateStock(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		delta   int
		repoErr error
		wantErr error
	}{
		{
			name:    "success",
			id:      "1",
			delta:   10,
			repoErr: nil,
			wantErr: nil,
		},
		{
			name:    "product not found",
			id:      "999",
			delta:   10,
			repoErr: ErrNotFound,
			wantErr: ErrNotFound,
		},
		{
			name:    "insufficient stock",
			id:      "1",
			delta:   -999,
			repoErr: ErrInsufficientStock,
			wantErr: ErrInsufficientStock,
		},
		{
			name:    "delta is zero",
			id:      "1",
			delta:   0,
			repoErr: nil,
			wantErr: ErrValidation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			repo := &mockRepo{
				updateStockFn: func(id string, delta int) error { return tt.repoErr },
			}
			svc := NewProductService(repo)

			err := svc.UpdateStock(tt.id, tt.delta)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("got err %v, want %v", err, tt.wantErr)
			}
		})
	}
}
