package usecase

import (
	"fmt"
	"unitrip/internal/entity"
)

type Admin interface {
	GetAllBuses() ([]*entity.Bus, error)
	GetBusByID(id int64) (*entity.Bus, error)
	CreateBus(bus *entity.Bus) (*entity.Bus, error)
	UpdateBus(bus *entity.Bus) error
	DeleteBus(id int64) error
}

type admin struct {
	busRepo BusRepo
}

func NewAdminService(busRepo BusRepo) Admin {
	return &admin{busRepo: busRepo}
}

func (u *admin) GetAllBuses() ([]*entity.Bus, error) {
	buses, err := u.busRepo.GetAll()
	if err != nil {
		return nil, err
	}
	if len(buses) == 0 {
		return nil, fmt.Errorf("%w:%v", ErrNotFound, "buses not found")
	}
	return buses, nil
}

func (u *admin) GetBusByID(id int64) (*entity.Bus, error) {
	bus, err := u.busRepo.GetByID(id)
	if err != nil {
		if err == ErrNotFound {
			return nil, fmt.Errorf("%w:%v", err, "bus not found")
		}
		return nil, err
	}
	return bus, nil
}

func (u *admin) CreateBus(bus *entity.Bus) (*entity.Bus, error) {
	bus, err := u.busRepo.Create(bus)
	if err != nil {
		return nil, err
	}
	return u.busRepo.GetByID(bus.ID)
}

func (u *admin) UpdateBus(bus *entity.Bus) error {
	return u.busRepo.Update(bus)
}

func (u *admin) DeleteBus(id int64) error {
	return u.busRepo.Delete(id)
}
