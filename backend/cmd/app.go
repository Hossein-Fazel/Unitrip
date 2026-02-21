package cmd

import (
	"errors"
	"log"
	"time"

	"unitrip/internal/adapter/http"
	"unitrip/internal/adapter/repository"
	"unitrip/internal/config"
	"unitrip/internal/infrastructure/database"
	"unitrip/internal/infrastructure/jwt"
	"unitrip/internal/infrastructure/web"
	"unitrip/internal/usecase"
)

func Run() error {
	// Config
	conf, err := config.New()
	if err != nil {
		return err
	}

	// Database
	db := database.NewDB(conf)
	if db == nil {
		log.Fatalln("db is nil ...")
	}

	// JWT
	jwtService := jwt.NewJWTService(conf.SecretKey, 24*time.Hour)

	// Repositories
	userRepo := repository.NewUserRepo(db)
	busRepo := repository.NewBusRepo(db)
	cityRepo := repository.NewCityRepo(db)

	// Services
	userService := usecase.NewUserService(userRepo)
	adminService := usecase.NewAdminService(busRepo, cityRepo)

	// Create admin account
	_, err = userService.CreateAdmin(conf.AdminUsername, conf.AdminPassword)
	if err != nil && !errors.Is(err, usecase.ErrUserAlreadyExist) {
		log.Fatal(err)
	}

	// Handlers
	userHandler := http.NewAuthHandler(userService, jwtService)
	adminHandler := http.NewAdminHandler(adminService)

	// Gin
	ginEngine := web.NewRouter()
	web.RegisterRoutes(ginEngine, userHandler, adminHandler, jwtService)

	// Run app
	log.Printf("Starting server on port %s", conf.WebPort)
	return ginEngine.Run(":" + conf.WebPort)
}
