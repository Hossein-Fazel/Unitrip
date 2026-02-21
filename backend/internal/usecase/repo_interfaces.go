package usecase

import (
	"unitrip/internal/entity"
)

type UserRepo interface {
	Save(user *entity.User) (*entity.User, error)
	Exist(user *entity.User) (bool, error)
	GetUserByEmail(email string) (*entity.User, error)
	GetUserByUsername(username string) (*entity.User, error)
}

type BusRepo interface {
	GetAll() ([]*entity.Bus, error)
	GetByID(id int64) (*entity.Bus, error)
	Create(bus *entity.Bus) (*entity.Bus, error)
	Update(bus *entity.Bus) error
	Delete(id int64) error
	GenerateSeats(id int64, n int) error
	GetSeatsByID(id int64) ([]*entity.Seat, error)
}

type CityRepo interface {
	Create(city *entity.City) (*entity.City, error)
	GetByName(name string) (*entity.City, error)
}