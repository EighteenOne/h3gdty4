package memory

import (
	"applicationDesignTest/internal/domain"
	"context"
	"errors"
	"sync"
)

type InMemoryOrderRepository struct {
	orders []domain.Order
	mu     sync.RWMutex
}

func NewInMemoryOrderRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{
		orders: make([]domain.Order, 0),
	}
}

func (r *InMemoryOrderRepository) Save(ctx context.Context, order *domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orders = append(r.orders, *order)
	return nil
}

func (r *InMemoryOrderRepository) CreateSnapshot(ctx context.Context) (interface{}, error) {
	if r.orders == nil {
		return nil, errors.New("repository is not initialized")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	snapshot := make([]domain.Order, len(r.orders))
	copy(snapshot, r.orders)
	return snapshot, nil
}

func (r *InMemoryOrderRepository) RestoreSnapshot(ctx context.Context, snapshot interface{}) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	snapshotData, ok := snapshot.([]domain.Order)
	if !ok {
		return errors.New("invalid snapshot type")
	}

	r.orders = snapshotData
	return nil
}
