package usecase

import (
	"unitrip/internal/entity"
)

type UserRepo interface {
	Save(user *entity.User) (*entity.User, error)
	Exist(user *entity.User) (bool, error)
}