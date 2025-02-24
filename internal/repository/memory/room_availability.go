package memory

import (
	"applicationDesignTest/internal/domain"
	"context"
	"errors"
	"sync"
	"time"
)

type InMemoryAvailabilityRepository struct {
	rooms map[string]domain.RoomAvailability
	mu    sync.RWMutex
}

func NewInMemoryAvailabilityRepository() *InMemoryAvailabilityRepository {
	return &InMemoryAvailabilityRepository{
		rooms: make(map[string]domain.RoomAvailability),
	}
}

func generateKey(hotelID, roomID string, date time.Time) string {
	return hotelID + "_" + roomID + "_" + date.Format("2006-01-02")
}

func (r *InMemoryAvailabilityRepository) GetAvailabilitiesForUpdate(ctx context.Context, hotelID, roomID string, from, to time.Time) ([]domain.RoomAvailability, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var availabilities []domain.RoomAvailability
	for d := from; !d.After(to); d = d.AddDate(0, 0, 1) {
		key := generateKey(hotelID, roomID, d)
		if availability, ok := r.rooms[key]; ok {
			availabilities = append(availabilities, availability)
		} else {
			// если не нашли вернем пустышку c датой
			availabilities = append(availabilities, domain.RoomAvailability{
				Date: d,
			})
		}
	}

	return availabilities, nil
}

func (r *InMemoryAvailabilityRepository) UpdateAvailabilities(ctx context.Context, availabilities []domain.RoomAvailability) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, availability := range availabilities {
		key := generateKey(availability.HotelID, availability.RoomID, availability.Date)
		r.rooms[key] = availability
	}
	return nil
}

func (r *InMemoryAvailabilityRepository) CreateSnapshot(ctx context.Context) (interface{}, error) {
	if r.rooms == nil {
		return nil, errors.New("repository is not initialized")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	snapshot := make(map[string]domain.RoomAvailability, len(r.rooms))
	for key, value := range r.rooms {
		snapshot[key] = value
	}

	return snapshot, nil
}

func (r *InMemoryAvailabilityRepository) RestoreSnapshot(ctx context.Context, snapshot interface{}) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	snapshotData, ok := snapshot.(map[string]domain.RoomAvailability)
	if !ok {
		return errors.New("invalid snapshot type")
	}

	r.rooms = snapshotData
	return nil
}
