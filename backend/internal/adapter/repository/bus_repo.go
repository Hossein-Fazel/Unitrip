package repository

import (
	"database/sql"
	"errors"
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
	query := `
		SELECT
			b.id,
			b.source_city_id,
			b.dest_city_id,
			b.travel_date,
			b.travel_time,
			b.price,
			c1.id, c1.name,
			c2.id, c2.name
		FROM buses b
		JOIN cities c1 ON c1.id = b.source_city_id
		JOIN cities c2 ON c2.id = b.dest_city_id
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

		var dateStr, timeStr string

		err := rows.Scan(
			&bus.ID,
			&source.ID,
			&dest.ID,
			&dateStr,
			&timeStr,
			&bus.Price,
			&source.ID, &source.Name,
			&dest.ID, &dest.Name,
		)
		if err != nil {
			return nil, err
		}

		dt, _ := time.Parse(
			"2006-01-02 15:04:05",
			dateStr+" "+timeStr,
		)
		bus.Schedule = dt

		bus.Route = entity.Route{
			Source: source,
			Dest:   dest,
		}

		buses = append(buses, &bus)
	}

	return buses, nil
}

func (b *bus) GetByID(id int64) (*entity.Bus, error) {
	query := `
		SELECT
			b.id,
			b.travel_date,
			b.travel_time,
			b.price,
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

	var dateStr, timeStr string

	err := b.db.QueryRow(query, id).Scan(
		&bus.ID,
		&dateStr,
		&timeStr,
		&bus.Price,
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

	dt, _ := time.Parse(
		"2006-01-02 15:04:05",
		dateStr+" "+timeStr,
	)
	bus.Schedule = dt

	bus.Route = entity.Route{
		Source: source,
		Dest:   dest,
	}

	return &bus, nil
}

func (b *bus) Create(bus *entity.Bus) (*entity.Bus, error) {
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

func (r *bus) Update(bus *entity.Bus) error {
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

	_, err := r.db.Exec(
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
	query := `DELETE FROM buses WHERE id = $1`
	_, err := b.db.Exec(query, id)
	return err
}
