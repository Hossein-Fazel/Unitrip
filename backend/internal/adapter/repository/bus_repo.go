package repository

import (
	"database/sql"
	"errors"
	"strings"
	"time"
	"unitrip/internal/entity"
	"unitrip/internal/usecase"
)

type bus struct {
	db *sql.DB
}

func NewBusRepo(db *sql.DB) usecase.BusRepo {
	return &bus{
		db: db,
	}
}

func (b *bus) GetAll() ([]*entity.Bus, error) {
	if b.db == nil {
		return nil, errors.New("database connection is nil")
	}

	query := `
		SELECT
			b.id,
			b.source_city_id,
			b.dest_city_id,
			b.travel_date,
			b.travel_time,
			b.price,
			b.rating_avg,
			b.rating_count,
			c1.id, c1.name,
			c2.id, c2.name
		FROM buses b
		JOIN cities c1 ON c1.id = b.source_city_id
		JOIN cities c2 ON c2.id = b.dest_city_id
		ORDER BY b.id DESC
	`

	rows, err := b.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var buses []*entity.Bus

	for rows.Next() {
		var bus entity.Bus
		var source entity.City
		var dest entity.City
		var rating entity.Rating

		var dateStr, timeStr string

		err := rows.Scan(
			&bus.ID,
			&source.ID,
			&dest.ID,
			&dateStr,
			&timeStr,
			&bus.Price,
			&rating.Average,
			&rating.Count,
			&source.ID, &source.Name,
			&dest.ID, &dest.Name,
		)
		if err != nil {
			return nil, err
		}

		dateStr = strings.Split(dateStr, "T")[0]
		timeStr = strings.TrimSuffix(strings.Split(timeStr, "T")[1], "Z")

		dt, _ := time.Parse(
			"2006-01-02 15:04:05",
			dateStr+" "+timeStr,
		)
		bus.Schedule = dt

		bus.Route = &entity.Route{
			Source: &source,
			Dest:   &dest,
		}
		bus.Rating = &rating

		buses = append(buses, &bus)
	}

	return buses, nil
}

func (b *bus) GetByID(id int64) (*entity.Bus, error) {
	if b.db == nil {
		return nil, errors.New("database connection is nil")
	}

	query := `
		SELECT
			b.id,
			b.travel_date,
			b.travel_time,
			b.price,
			b.rating_avg,
			b.rating_count,
			c1.id, c1.name,
			c2.id, c2.name
		FROM buses b
		JOIN cities c1 ON c1.id = b.source_city_id
		JOIN cities c2 ON c2.id = b.dest_city_id
		WHERE b.id = $1
	`

	var bus entity.Bus
	var source entity.City
	var dest entity.City
	var rating entity.Rating

	var dateStr, timeStr string

	err := b.db.QueryRow(query, id).Scan(
		&bus.ID,
		&dateStr,
		&timeStr,
		&bus.Price,
		&rating.Average,
		&rating.Count,
		&source.ID, &source.Name,
		&dest.ID, &dest.Name,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, usecase.ErrNotFound
		default:
			return nil, err
		}
	}

	dateStr = strings.Split(dateStr, "T")[0]
	timeStr = strings.TrimSuffix(strings.Split(timeStr, "T")[1], "Z")

	dt, _ := time.Parse(
		"2006-01-02 15:04:05",
		dateStr+" "+timeStr,
	)
	bus.Schedule = dt

	bus.Route = &entity.Route{
		Source: &source,
		Dest:   &dest,
	}
	bus.Rating = &rating

	return &bus, nil
}

func (b *bus) Create(bus *entity.Bus) (*entity.Bus, error) {
	if b.db == nil {
		return nil, errors.New("database connection is nil")
	}

	query := `
		INSERT INTO buses (
			source_city_id,
			dest_city_id,
			travel_date,
			travel_time,
			price
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	date := bus.Schedule.Format("2006-01-02")
	time := bus.Schedule.Format("15:04:05")

	err := b.db.QueryRow(
		query,
		bus.Route.Source.ID,
		bus.Route.Dest.ID,
		date,
		time,
		bus.Price,
	).Scan(&bus.ID)

	if err != nil {
		return nil, err
	}

	return bus, nil
}

func (b *bus) Update(bus *entity.Bus) error {
	if b.db == nil {
		return errors.New("database connection is nil")
	}

	query := `
		UPDATE buses
		SET
			source_city_id = $1,
			dest_city_id = $2,
			travel_date = $3,
			travel_time = $4,
			price = $5
		WHERE id = $6
	`

	date := bus.Schedule.Format("2006-01-02")
	time := bus.Schedule.Format("15:04:05")

	_, err := b.db.Exec(
		query,
		bus.Route.Source.ID,
		bus.Route.Dest.ID,
		date,
		time,
		bus.Price,
		bus.ID,
	)

	return err
}

func (b *bus) Delete(id int64) error {
	if b.db == nil {
		return errors.New("database connection is nil")
	}

	query := `DELETE FROM buses WHERE id = $1`
	_, err := b.db.Exec(query, id)
	return err
}

func (b *bus) GenerateSeats(id int64, n int) error {
	if b.db == nil {
		return errors.New("database connection is nil")
	}

	query := `
		INSERT INTO bus_seats (bus_id, seat_no)
		SELECT $1, generate_series(1, $2)
	`
	_, err := b.db.Exec(query, id, n)
	return err
}

func (b *bus) GetSeatsByID(id int64) ([]*entity.Seat, error) {
	if b.db == nil {
		return nil, errors.New("database connection is nil")
	}

	query := `
		SELECT id, seat_no, status
		FROM bus_seats
		WHERE bus_id = $1
		ORDER BY seat_no
	`

	rows, err := b.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seats []*entity.Seat
	for rows.Next() {
		seat := &entity.Seat{}
		if err := rows.Scan(
			&seat.ID,
			&seat.SeatNo,
			&seat.Status,
		); err != nil {
			return nil, err
		}
		seats = append(seats, seat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return seats, nil
}
