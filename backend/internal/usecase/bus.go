package usecase

import (
	"errors"
	"fmt"
	"unitrip/internal/entity"
)

type Bus interface {
	List() ([]*entity.Bus, error)
	GetByID(id int64) (*entity.Bus, error)
	Create(bus *entity.Bus) (*entity.Bus, error)
	Update(bus *entity.Bus) error
	Delete(id int64) error
}

type bus struct {
	busRepo  BusRepo
	cityRepo CityRepo
}

func NewBusService(busRepo BusRepo, cityRepo CityRepo) Bus {
	return &bus{
		busRepo:  busRepo,
		cityRepo: cityRepo,
	}
}

func (b *bus) List() ([]*entity.Bus, error) {
	buses, err := b.busRepo.GetAll()
	if err != nil {
		return nil, err
	}
	if len(buses) == 0 {
		return nil, fmt.Errorf("%w:%v", ErrNotFound, "buses not found")
	}

	for _, bus := range buses {
		seats, err := b.busRepo.GetSeatsByID(bus.ID)
		if err != nil {
			return nil, err
		}

		bus.Seats = seats
	}
	return buses, nil
}

func (b *bus) GetByID(id int64) (*entity.Bus, error) {
	bus, err := b.busRepo.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			return nil, fmt.Errorf("%w:%v", err, "bus not found")
		default:
			return nil, err
		}
	}

	seats, err := b.busRepo.GetSeatsByID(bus.ID)
	if err != nil {
		return nil, err
	}

	bus.Seats = seats

	return bus, nil
}

func (b *bus) Create(bus *entity.Bus) (*entity.Bus, error) {
	var err error
	bus.Route.Source, err = b.cityRepo.GetByName(bus.Route.Source.Name)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			return nil, fmt.Errorf("%w:%v", err, "no city with this name")
		default:
			return nil, fmt.Errorf("%w:%v", ErrInternal, err.Error())
		}
	}

	bus.Route.Dest, err = b.cityRepo.GetByName(bus.Route.Dest.Name)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			return nil, fmt.Errorf("%w:%v", err, "no city with this name")
		default:
			return nil, fmt.Errorf("%w:%v", ErrInternal, err.Error())
		}
	}

	bus, err = b.busRepo.Create(bus)
	if err != nil {
		return nil, err
	}

	err = b.busRepo.GenerateSeats(bus.ID, 25)
	if err != nil {
		return nil, err
	}

	return bus, nil
}

func (b *bus) Update(bus *entity.Bus) error {
	var err error
	bus.Route.Source, err = b.cityRepo.GetByName(bus.Route.Source.Name)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			return fmt.Errorf("%w:%v", err, "no city with this name")
		default:
			return fmt.Errorf("%w:%v", ErrInternal, err.Error())
		}
	}

	bus.Route.Dest, err = b.cityRepo.GetByName(bus.Route.Dest.Name)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			return fmt.Errorf("%w:%v", err, "no city with this name")
		default:
			return fmt.Errorf("%w:%v", ErrInternal, err.Error())
		}
	}
	return b.busRepo.Update(bus)
}

func (b *bus) Delete(id int64) error {
	return b.busRepo.Delete(id)
}
