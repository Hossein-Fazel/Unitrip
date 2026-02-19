package http

import (
	"unitrip/internal/entity"
	"unitrip/internal/infrastructure/jwt"
	"unitrip/internal/usecase"

	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	Signup(c *gin.Context)
}

type user struct {
	authService usecase.User
	jwtService  jwt.JWTService
}

func NewUserHandler(authService usecase.User, jwtService jwt.JWTService) UserHandler {
	return &user{
		authService: authService,
		jwtService:  jwtService,
	}
}

func (a *user) Signup(c *gin.Context) {
	var signupRequest SignupRequest

	if err := c.ShouldBindJSON(&signupRequest); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request"})
		return
	}

	user := &entity.User{
		Username: signupRequest.Username,
		Email:    signupRequest.Email,
		Password: signupRequest.Password,
	}
	user, err := a.authService.Signup(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// generate token for this user
	jwtToken, err := a.jwtService.GenerateToken(user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "could not generate token"})
		return
	}

	c.JSON(http.StatusCreated, SignUpResponse{
		Message: "User registered successfully",
		Token:   jwtToken,
	})
}
