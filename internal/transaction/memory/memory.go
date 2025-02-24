package memory

import (
	"context"
	"errors"
	"sync"
)

type Shapshoter interface {
	CreateSnapshot(ctx context.Context) (interface{}, error)
	RestoreSnapshot(ctx context.Context, snapshot interface{}) error
}

type InMemoryTransactionManager struct {
	mu           sync.Mutex
	participants []Shapshoter
}

func NewInMemoryTransactionManager(participants ...Shapshoter) *InMemoryTransactionManager {
	return &InMemoryTransactionManager{
		participants: participants,
	}
}

func (tm *InMemoryTransactionManager) RunInTransaction(ctx context.Context, body func(ctx context.Context) error) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	snapshots := make([]interface{}, len(tm.participants))
	for i, participant := range tm.participants {
		snapshot, err := participant.CreateSnapshot(ctx)
		if err != nil {
			return err
		}
		snapshots[i] = snapshot
	}

	err := body(ctx)
	if err != nil {
		for i, participant := range tm.participants {
			if err := participant.RestoreSnapshot(ctx, snapshots[i]); err != nil {
				return errors.New("failed to restore snapshot: " + err.Error())
			}
		}
		return err
	}

	snapshots = make([]interface{}, 0)
	return nil
}
