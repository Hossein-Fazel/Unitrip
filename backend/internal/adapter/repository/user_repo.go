package repository

import (
	"database/sql"
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
