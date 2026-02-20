package usecase

import (
	"fmt"
	"unitrip/internal/entity"

	"github.com/tailscale/golang-x-crypto/bcrypt"
)

type User interface {
	Signup(user *entity.User) (*entity.User, error)
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
	user.Password = hash_password(user.Password)
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

func hash_password(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	return  string(hashedPassword)
}
