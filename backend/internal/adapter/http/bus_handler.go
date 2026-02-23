package http

import (
	"errors"
	"net/http"
	"strconv"
	"time"
	"unitrip/internal/entity"
	"unitrip/internal/usecase"

	"github.com/gin-gonic/gin"
)

type BusHandler interface {
	List(c *gin.Context)
	GetByID(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type bus struct {
	busService usecase.Bus
}

func NewBusHandler(busService usecase.Bus) BusHandler {
	return &bus{busService: busService}
}

func (h *bus) List(c *gin.Context) {
	buses, err := h.busService.List()
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "no buses found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, generateBusesResponse(buses))
}

func (h *bus) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "invalid parameter"})
		return
	}

	bus, err := h.busService.GetByID(id)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "bus not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, generateBusResponse(bus))
}

func (h *bus) Create(c *gin.Context) {
	var request BusRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	inputBus, err := generateBus(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	_, err = h.busService.Create(inputBus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, MessageResponse{
		Message: "bus create successfully",
	})
}

func (h *bus) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "invalid parameter"})
		return
	}

	var request BusRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	input, err := generateBus(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	input.ID = id
	if err := h.busService.Update(input); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{
		Message: "bus update successfully",
	})
}

func (h *bus) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "invalid parameter"})
		return
	}

	if err := h.busService.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, MessageResponse{
		Message: "bus deleted successfully",
	})
}

func generateBusesResponse(buses []*entity.Bus) []*BusResponse {
	var response []*BusResponse
	for _, bus := range buses {
		var busResponse BusResponse
		busResponse.ID = bus.ID
		busResponse.SourceCity = bus.Route.Source.Name
		busResponse.DestCity = bus.Route.Dest.Name
		busResponse.Price = bus.Price
		busResponse.TravelDate = bus.Schedule.Format("2006-01-02")
		busResponse.TravelTime = bus.Schedule.Format("15:04:05")
		busResponse.Rating.Average = bus.Rating.Average
		busResponse.Rating.Count = bus.Rating.Count
		busResponse.AvailableSeatsCount = bus.AvailableSeatsCount()

		response = append(response, &busResponse)
	}

	return response
}

func generateBusResponse(bus *entity.Bus) DetailedBusResponse {
	var busResponse DetailedBusResponse

	busResponse.ID = bus.ID
	busResponse.SourceCity = bus.Route.Source.Name
	busResponse.DestCity = bus.Route.Dest.Name
	busResponse.Price = bus.Price
	busResponse.TravelDate = bus.Schedule.Format("2006-01-02")
	busResponse.TravelTime = bus.Schedule.Format("15:04:05")
	busResponse.Rating.Average = bus.Rating.Average
	busResponse.Rating.Count = bus.Rating.Count

	for _, seat := range bus.Seats {
		busResponse.Seats = append(busResponse.Seats, Seat{
			ID:     seat.ID,
			SeatNo: seat.SeatNo,
			Status: seat.Status,
		})
	}

	return busResponse
}

func generateBus(busRequest BusRequest) (*entity.Bus, error) {
	var bus entity.Bus

	route := entity.Route{
		Source: &entity.City{Name: busRequest.SourceCity},
		Dest:   &entity.City{Name: busRequest.DestCity},
	}
	bus.Route = &route
	bus.Price = busRequest.Price

	layout := "2006-01-02 15:04"
	dt, err := time.Parse(layout, busRequest.Date+" "+busRequest.Time)
	if err != nil {
		return nil, err
	}
	bus.Schedule = dt

	return &bus, nil
}
