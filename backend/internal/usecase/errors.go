package usecase

import "errors"

var (
	ErrInternal         = errors.New("internal error accured")
	ErrNotFound         = errors.New("not found")
	ErrUserAlreadyExist = errors.New("user already exist")
	ErrUserSaveFailed   = errors.New("failed to save user")
	ErrInvalidRequest   = errors.New("Invalid request")
	ErrPasswordWrong    = errors.New("password wrong")
)
