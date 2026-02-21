package http

import (
	"errors"
	"regexp"
	"unitrip/internal/entity"
	"unitrip/internal/usecase"

	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	Signup(c *gin.Context)
	Login(c *gin.Context)
}

type auth struct {
	userService usecase.User
	jwtService  JWTService
}

func NewAuthHandler(userService usecase.User, jwtService JWTService) AuthHandler {
	return &auth{
		userService: userService,
		jwtService:  jwtService,
	}
}

func (a *auth) Signup(c *gin.Context) {
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
	user, err := a.userService.Signup(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// generate token for this user
	jwtToken, err := a.jwtService.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "could not generate token"})
		return
	}

	c.JSON(http.StatusCreated, SignUpResponse{
		Message: "register successful",
		Token:   jwtToken,
	})
}

func (a *auth) Login(c *gin.Context) {
	var loginRequest LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request"})
		return
	}

	user := &entity.User{
		Password: loginRequest.Password,
	}
	if isEmail(loginRequest.Identifier) {
		user.Email = loginRequest.Identifier
	} else {
		user.Username = loginRequest.Identifier
	}

	user, err := a.userService.Login(user)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInternal) :
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: usecase.ErrInternal.Error(),
			})

		case errors.Is(err, usecase.ErrPasswordWrong), errors.Is(err, usecase.ErrNotFound) :
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error: "invalid credentials",
			})
		}
		return
	}

	// generate token for this user
	jwtToken, err := a.jwtService.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "could not generate token"})
		return
	}

	c.JSON(http.StatusCreated, SignUpResponse{
		Message: "login successful",
		Token:   jwtToken,
	})
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func isEmail(s string) bool {
	return emailRegex.MatchString(s)
}
