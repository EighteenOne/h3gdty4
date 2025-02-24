package domain

import (
	"errors"
	"time"
)

type Order struct {
	HotelID   string
	RoomID    string
	UserEmail string
	From      time.Time
	To        time.Time
}

func NewOrder(hotelID string, roomID string, userEmail string, from, to time.Time) (*Order, error) {
	if hotelID == "" {
		return nil, errors.New("hotel ID is required")
	}
	if roomID == "" {
		return nil, errors.New("room ID is required")
	}
	if userEmail == "" {
		return nil, errors.New("email is required")
	}
	if from.After(to) {
		return nil, errors.New("'from' date cannot be after 'to' date")
	}

	return &Order{
		HotelID:   hotelID,
		RoomID:    roomID,
		UserEmail: userEmail,
		From:      from,
		To:        to,
	}, nil
}
