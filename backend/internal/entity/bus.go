package entity

import "time"

type Bus struct {
	ID        int64
	Route     Route
	Schedule  Schedule
	Price     int64
	Seats     []Seat
	Rating    Rating
	CreatedAt time.Time
}

type Route struct {
	Source City
	Dest   City
}

type Schedule struct {
	TravelDate time.Time
	TravelTime time.Time
}

type Seat struct {
	SeatNo int
	Status SeatStatus
}

type SeatStatus string

const (
	SeatAvailable SeatStatus = "AVAILABLE"
	SeatPending   SeatStatus = "PENDING"
	SeatBooked    SeatStatus = "BOOKED"
)

func (b Bus) AvailableSeatsCount() int {
	count := 0
	for _, seat := range b.Seats {
		if seat.Status == SeatAvailable {
			count++
		}
	}
	return count
}

func (b Bus) IsAvailable() bool {
	return time.Now().Before(b.Schedule.TravelDate)
}