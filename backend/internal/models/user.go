package models

import "time"

type UserRole string

const (
	RoleAdmin UserRole = "ADMIN"
	RoleUser  UserRole = "USER"
)

type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"passwordHash"`
	Role         UserRole  `json:"role"`
	CreatedAt    time.Time `json:"createdAt"`
}
