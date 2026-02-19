package cmd

import (
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
	conf, err := config.New()
	if err != nil {
		return err
	}

	db := database.NewDB(conf)
	if db == nil {
		log.Fatalln("db is nil ...")
	}

	jwtService := jwt.NewJWTService(conf.SecretKey, 24*time.Hour)

	userRepo := repository.NewUserRepo(db)
	userService := usecase.NewUserService(userRepo)
	userHandler := http.NewAuthHandler(userService, jwtService)

	ginEngine := web.NewRouter()
	web.RegisterRoutes(ginEngine, userHandler)

	log.Printf("Starting server on port %s", conf.WebPort)
	return ginEngine.Run(":" + conf.WebPort)
}
