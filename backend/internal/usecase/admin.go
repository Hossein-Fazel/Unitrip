package usecase

import (
	"errors"
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
	busRepo  BusRepo
	cityRepo CityRepo
}

func NewAdminService(busRepo BusRepo, cityRepo CityRepo) Admin {
	return &admin{
		busRepo: busRepo,
		cityRepo: cityRepo,
	}
}

func (a *admin) GetAllBuses() ([]*entity.Bus, error) {
	buses, err := a.busRepo.GetAll()
	if err != nil {
		return nil, err
	}
	if len(buses) == 0 {
		return nil, fmt.Errorf("%w:%v", ErrNotFound, "buses not found")
	}

	for _, bus := range(buses) {
		seats, err := a.busRepo.GetSeatsByID(bus.ID)
		if err != nil {
			return nil, err
		}

		bus.Seats = seats
	}
	return buses, nil
}

func (a *admin) GetBusByID(id int64) (*entity.Bus, error) {
	bus, err := a.busRepo.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			return nil, fmt.Errorf("%w:%v", err, "bus not found")
		default:
			return nil, err
		}
	}

	seats, err := a.busRepo.GetSeatsByID(bus.ID)
	if err != nil {
		return nil, err
	}

	bus.Seats = seats
	
	return bus, nil
}

func (a *admin) CreateBus(bus *entity.Bus) (*entity.Bus, error) {
	var err error
	bus.Route.Source, err = a.cityRepo.GetByName(bus.Route.Source.Name)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			return nil, fmt.Errorf("%w:%v", err, "no city with this name")
		default:
			return nil, fmt.Errorf("%w:%v", ErrInternal, err.Error())
		}
	}

	bus.Route.Dest, err = a.cityRepo.GetByName(bus.Route.Dest.Name)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			return nil, fmt.Errorf("%w:%v", err, "no city with this name")
		default:
			return nil, fmt.Errorf("%w:%v", ErrInternal, err.Error())
		}
	}

	bus, err = a.busRepo.Create(bus)
	if err != nil {
		return nil, err
	}

	err = a.busRepo.GenerateSeats(bus.ID, 25)
	if err != nil {
		return nil, err
	}

	return bus, nil
}

func (a *admin) UpdateBus(bus *entity.Bus) error {
	var err error
	bus.Route.Source, err = a.cityRepo.GetByName(bus.Route.Source.Name)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			return fmt.Errorf("%w:%v", err, "no city with this name")
		default:
			return fmt.Errorf("%w:%v", ErrInternal, err.Error())
		}
	}

	bus.Route.Dest, err = a.cityRepo.GetByName(bus.Route.Dest.Name)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			return fmt.Errorf("%w:%v", err, "no city with this name")
		default:
			return fmt.Errorf("%w:%v", ErrInternal, err.Error())
		}
	}
	return a.busRepo.Update(bus)
}

func (a *admin) DeleteBus(id int64) error {
	return a.busRepo.Delete(id)
}
