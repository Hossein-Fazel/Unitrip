package usecase

import (
	"errors"
	"fmt"
	"unitrip/internal/entity"

	"github.com/tailscale/golang-x-crypto/bcrypt"
)

type User interface {
	Signup(user *entity.User) (*entity.User, error)
    Login(user *entity.User) (*entity.User, error)
}
type user struct {
	userRepo UserRepo
}

func NewUserService(userRepo UserRepo) User {
	return &user{
		userRepo: userRepo,
	}
}

func (a *user) Signup(user *entity.User) (*entity.User, error) {
	user.Password = hashPassword(user.Password)
	user.Role = "USER"

	exist, err := a.userRepo.Exist(user)
	if err != nil {
		return nil, fmt.Errorf("%w:%v", ErrInternal, err)
	}
	if exist {
		return nil, ErrUserAlreadyExist
	}

	user, err = a.userRepo.Save(user)
	if err != nil {
		return nil, fmt.Errorf("%x:%v", ErrUserSaveFailed, err)
	}

	return user, nil
}

func (a *user) Login(inputUser *entity.User) (*entity.User, error) {
	var gotUser *entity.User
	var err error

	if inputUser.Username != "" {
		fmt.Println("login with username")
		gotUser, err = a.userRepo.GetUserByUsername(inputUser.Username)
	} else if inputUser.Email != "" {
		fmt.Println("login with email")

		gotUser, err = a.userRepo.GetUserByEmail(inputUser.Email)
	} else {
		return nil, fmt.Errorf("%w:%v", ErrInvalidRequest, "username or email is required for login")
	}

	if err != nil {
		switch {
		case errors.Is(err, ErrUserNotFound) :
			return nil, err
		default:
			return nil, fmt.Errorf("%w:%v", ErrInternal, err)
		}
	}

	if !checkPassword(gotUser.Password, inputUser.Password){
		return nil, ErrPasswordWrong
	}

	return gotUser, nil
}

func hashPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	return  string(hashedPassword)
}

func checkPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)
	return err == nil
}