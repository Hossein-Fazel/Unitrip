package repository

import (
	"database/sql"
	"errors"
	"unitrip/internal/entity"
	"unitrip/internal/usecase"
)

type user struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) usecase.UserRepo {
	return &user{
		db: db,
	}
}

func (u *user) Save(user *entity.User) (*entity.User, error) {
	if u.db == nil {
		return nil, errors.New("database connection is nil")
	}

	query := `
		INSERT INTO users (username, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err := u.db.QueryRow(
		query,
		user.Username,
		user.Email,
		user.Password,
		user.Role,
	).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *user) Exist(user *entity.User) (bool, error) {
	if u.db == nil {
		return false, errors.New("database connection is nil")
	}

	var exists bool

	query := `
		SELECT EXISTS (
			SELECT 1 FROM users
			WHERE username = $1 OR email = $2
		)
	`

	err := u.db.QueryRow(query, user.Username, user.Email).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (u *user) GetUserByEmail(email string) (*entity.User, error) {
	if u.db == nil {
		return nil, errors.New("database connection is nil")
	}

	query := `
		SELECT id, username, email, password_hash, role, created_at
		FROM users
		WHERE email = $1
	`
	var user entity.User

	err := u.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, usecase.ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil

}

func (u *user) GetUserByUsername(username string) (*entity.User, error) {
	if u.db == nil {
		return nil, errors.New("database connection is nil")
	}

	query := `
		SELECT id, username, email, password_hash, role, created_at
		FROM users
		WHERE username = $1
	`
	var user entity.User

	err := u.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, usecase.ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
