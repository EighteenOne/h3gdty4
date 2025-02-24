package service

import (
	"applicationDesignTest/internal/domain"
	log "applicationDesignTest/internal/logger"
	"context"
	"fmt"
	"time"
)

type orderRepository interface {
	Save(ctx context.Context, order *domain.Order) error
}

type availabilityRepository interface {
	GetAvailabilities(ctx context.Context, hotelID, roomID string, from, to time.Time) ([]domain.RoomAvailability, error)
	UpdateAvailabilities(ctx context.Context, availabilities []domain.RoomAvailability) error
}

type transactionManager interface {
	RunInTransaction(ctx context.Context, body func(txCtx context.Context) error) error
}

type OrderService struct {
	orderRepo        orderRepository
	availabilityRepo availabilityRepository

	txManager transactionManager
}

func NewOrderService(orderRepo orderRepository, availabilityRepo availabilityRepository, txManager transactionManager) *OrderService {
	return &OrderService{
		orderRepo:        orderRepo,
		availabilityRepo: availabilityRepo,
		txManager:        txManager,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, order *domain.Order) error {
	return s.txManager.RunInTransaction(ctx, func(txCtx context.Context) error {

		availabilities, err := s.availabilityRepo.GetAvailabilities(
			txCtx,
			order.HotelID,
			order.RoomID,
			order.From,
			order.To,
		)
		if err != nil {
			return fmt.Errorf("failed to get availability: %w", err)
		}

		unavailableDays := make([]time.Time, 0)
		availabilityUpdates := make([]domain.RoomAvailability, 0, len(availabilities))

		for _, availability := range availabilities {
			if availability.Quota < 1 {
				unavailableDays = append(unavailableDays, availability.Date)
				continue
			}

			availability.Quota--
			availabilityUpdates = append(availabilityUpdates, availability)
		}

		if len(unavailableDays) > 0 {
			log.Errorf("Hotel room is not available for selected dates:\n%v\n%v", order, unavailableDays)
			return domain.ErrRoomIsNotAvailable
		}

		if err := s.availabilityRepo.UpdateAvailabilities(txCtx, availabilityUpdates); err != nil {
			return fmt.Errorf("failed to update availability: %w", err)
		}

		if err := s.orderRepo.Save(txCtx, order); err != nil {
			return fmt.Errorf("failed to save order: %w", err)
		}

		return nil
	})
}
