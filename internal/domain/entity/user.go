package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserID uuid.UUID

func (id UserID) String() string {
	return id.UUID().String()
}

func (id UserID) UUID() uuid.UUID {
	return uuid.UUID(id)
}

type User struct {
	ID           UserID
	Login        string
	PasswordHash string
	DisplayName  string
}

type UserCreate struct {
	ID           UserID
	Login        string
	PasswordHash string
	DisplayName  string
}

type UserCreateRequest struct {
	Login       string
	Password    string
	DisplayName string
}

type AuthRequest struct {
	Login    string
	Password string
}

type AuthResult struct {
	AccessToken string
	ExpiresAt   time.Time
	User        User
}
