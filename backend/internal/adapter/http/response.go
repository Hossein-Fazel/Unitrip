package http

import (
	"unitrip/internal/entity"
)

// SignUpResponse is the response returned after successful signup
type SignUpResponse struct {
	Message string `json:"message" example:"User registered successfully"`
	Token   string `json:"token" example:"generated token"`
}

// ErrorResponse is returned when an error occurs
type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}

// MessageResponse is returned when successful and message occurs
type MessageResponse struct {
	Message string `json:"message" example:"this is a message"`
}

type BusResponse struct {
	ID                  int64          `json:"id"`
	SourceCity          string         `json:"source_city"`
	DestCity            string         `json:"dest_city"`
	AvailableSeatsCount int            `json:"available_seat_count"`
	TravelDate          string         `json:"travel_date"`
	TravelTime          string         `json:"travel_time"`
	Price               float64        `json:"price"`
	Rating              RatingResponse `json:"rating"`
}

type DetailedBusResponse struct {
	ID                  int64          `json:"id"`
	SourceCity          string         `json:"source_city"`
	DestCity            string         `json:"dest_city"`
	Seats               []Seat         `json:"seats"`
	TravelDate          string         `json:"travel_date"`
	TravelTime          string         `json:"travel_time"`
	Price               float64        `json:"price"`
	Rating              RatingResponse `json:"rating"`
}

type RatingResponse struct {
	Average float32 `json:"average"`
	Count   int     `json:"count"`
}

type Seat struct {
	SeatNo int               `json:"number"`
	Status entity.SeatStatus `json:"status"`
}
