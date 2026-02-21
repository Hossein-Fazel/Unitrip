package repository

import (
	"database/sql"
	"errors"
	"unitrip/internal/entity"
	"unitrip/internal/usecase"
)

type city struct {
	db *sql.DB
}

func NewCityRepo(db *sql.DB) usecase.CityRepo {
	return &city{
		db: db,
	}
}

func (c *city) Create(city *entity.City) (*entity.City, error) {
	if c.db == nil {
		return nil, errors.New("database connection is nil")
	}

	query := `
		INSERT INTO cities (name)
		VALUES ($1)
		RETURNING id
	`

	err := c.db.QueryRow(query, city.Name).Scan(&city.ID)
	if err != nil {
		return nil, err
	}

	return city, nil
}

func (c *city) GetByName(name string) (*entity.City, error) {
	if c.db == nil {
		return nil, errors.New("database connection is nil")
	}

	query := `
		SELECT id, name
		FROM cities
		WHERE name = $1
	`

	city := &entity.City{}
	err := c.db.QueryRow(query, name).Scan(
		&city.ID,
		&city.Name,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, usecase.ErrNotFound
		default:
			return nil, err
		}
	}

	return city, nil
}
