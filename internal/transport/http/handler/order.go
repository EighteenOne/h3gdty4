package handler

import (
	"applicationDesignTest/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	log "applicationDesignTest/internal/logger"
)

type request struct {
	HotelID   string    `json:"hotel_id"`
	RoomID    string    `json:"room_id"`
	UserEmail string    `json:"email"`
	From      time.Time `json:"from"`
	To        time.Time `json:"to"`
}

type response struct {
	HotelID   string    `json:"hotel_id"`
	RoomID    string    `json:"room_id"`
	UserEmail string    `json:"email"`
	From      time.Time `json:"from"`
	To        time.Time `json:"to"`
}

type OrderCreator interface {
	CreateOrder(ctx context.Context, order *domain.Order) error
}

func NewCreateOrderHandler(orderCreator OrderCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		newOrder, err := domain.NewOrder(req.HotelID, req.RoomID, req.UserEmail, req.From, req.To)
		if err != nil {
			responseError(w, fmt.Sprintf("Cannot create order: %v", err), http.StatusBadRequest, *newOrder)
		}

		err = orderCreator.CreateOrder(r.Context(), newOrder)
		if errors.Is(err, domain.ErrRoomIsNotAvailable) {
			responseError(w, "Hotel room is not available for selected dates", http.StatusBadRequest, *newOrder)
			return
		}
		if err != nil {
			responseError(w, "Something went wrong", http.StatusInternalServerError, *newOrder)
			return
		}

		responseOK(w, r, response{
			HotelID:   newOrder.HotelID,
			RoomID:    newOrder.RoomID,
			UserEmail: newOrder.UserEmail,
			From:      newOrder.From,
			To:        newOrder.To,
		})

		log.Info("Order successfully created: %v", newOrder)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, resp response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func responseError(w http.ResponseWriter, msg string, code int, order domain.Order) {
	http.Error(w, msg, code)
	log.Errorf("%s: %v", msg, order)
}
