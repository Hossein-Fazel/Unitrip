package http

import (
	"errors"
	"net/http"
	"strconv"
	"unitrip/internal/entity"
	"unitrip/internal/usecase"

	"github.com/gin-gonic/gin"
)

type AdminHandler interface {
	GetAllBuses(c *gin.Context)
	GetBusByID(c *gin.Context)
	CreateBus(c *gin.Context)
	UpdateBus(c *gin.Context)
	DeleteBus(c *gin.Context)
}

type admin struct {
	adminService usecase.Admin
}

func NewAdminHandler(adminService usecase.Admin) AdminHandler {
	return &admin{adminService: adminService}
}

func (h *admin) GetAllBuses(c *gin.Context) {
	buses, err := h.adminService.GetAllBuses()
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "no buses found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, buses)
}

func (h *admin) GetBusByID(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.ParseInt(idParam, 10, 64)
	bus, err := h.adminService.GetBusByID(id)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "bus not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bus)
}

func (h *admin) CreateBus(c *gin.Context) {
	var input entity.Bus
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := h.adminService.CreateBus(&input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, MessageResponse{
		Message: "bus create successfully",
	})
}

func (h *admin) UpdateBus(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.ParseInt(idParam, 10, 64)
	var input entity.Bus
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.ID = id
	if err := h.adminService.UpdateBus(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, MessageResponse{
		Message: "bus update successfully",
	})
}

func (h *admin) DeleteBus(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.ParseInt(idParam, 10, 64)
	if err := h.adminService.DeleteBus(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, MessageResponse{
		Message: "bus deleted successfully",
	})
}
