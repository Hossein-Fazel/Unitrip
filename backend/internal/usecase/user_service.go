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
	CreateAdmin(username, password string) (*entity.User, error)
}
type user struct {
	userRepo UserRepo
}

func NewUserService(userRepo UserRepo) User {
	return &user{
		userRepo: userRepo,
	}
}

func (u *user) Signup(user *entity.User) (*entity.User, error) {
	user.Password = hashPassword(user.Password)
	user.Role = entity.RoleUser

	exist, err := u.userRepo.Exist(user)
	if err != nil {
		return nil, fmt.Errorf("%w:%v", ErrInternal, err)
	}
	if exist {
		return nil, ErrUserAlreadyExist
	}

	user, err = u.userRepo.Save(user)
	if err != nil {
		return nil, fmt.Errorf("%x:%v", ErrUserSaveFailed, err)
	}

	return user, nil
}

func (u *user) Login(inputUser *entity.User) (*entity.User, error) {
	var gotUser *entity.User
	var err error

	if inputUser.Username != "" {
		gotUser, err = u.userRepo.GetUserByUsername(inputUser.Username)
	} else if inputUser.Email != "" {
		gotUser, err = u.userRepo.GetUserByEmail(inputUser.Email)
	} else {
		return nil, fmt.Errorf("%w:%v", ErrInvalidRequest, "username or email is required for login")
	}

	if err != nil {
		switch {
		case errors.Is(err, ErrUserNotFound):
			return nil, err
		default:
			return nil, fmt.Errorf("%w:%v", ErrInternal, err)
		}
	}

	if !checkPassword(gotUser.Password, inputUser.Password) {
		return nil, ErrPasswordWrong
	}

	return gotUser, nil
}

func (u *user) CreateAdmin(username, password string) (*entity.User, error) {
	adminUser := &entity.User{
		Username: username,
		Password: hashPassword(password),
		Role:     entity.RoleAdmin,
	}

	exist, err := u.userRepo.Exist(adminUser)
	if err != nil {
		return nil, fmt.Errorf("%w:%v", ErrInternal, err)
	}
	if exist {
		return nil, ErrUserAlreadyExist
	}

	adminUser, err = u.userRepo.Save(adminUser)
	if err != nil {
		return nil, fmt.Errorf("%x:%v", ErrUserSaveFailed, err)
	}

	return adminUser, nil
}

func hashPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	return string(hashedPassword)
}

func checkPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)
	return err == nil
}
