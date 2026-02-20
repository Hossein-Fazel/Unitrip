package entity

import "time"

type User struct {
	ID        int64
	Username  string
	Password  string
	Email     string
	Role      UserRole
	CreatedAt time.Time
}


type UserRole string

const (
	RoleUser  UserRole = "USER"
	RoleAdmin UserRole = "ADMIN"
)