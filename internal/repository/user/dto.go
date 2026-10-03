package user

import (
	"familyquest-backend/internal/domain/entity"

	"github.com/google/uuid"
)

const (
	table = "users"

	fieldID           = "id"
	fieldLogin        = "login"
	fieldPasswordHash = "password_hash"
	fieldDisplayName  = "display_name"
)

type user struct {
	ID          uuid.UUID `db:"id"`
	Login       string    `db:"login"`
	DisplayName string    `db:"display_name"`
}

func (u user) toEntity() entity.User {
	return entity.User{
		ID:          entity.UserID(u.ID),
		Login:       u.Login,
		DisplayName: u.DisplayName,
	}
}

type userCredentials struct {
	ID           uuid.UUID `db:"id"`
	Login        string    `db:"login"`
	PasswordHash string    `db:"password_hash"`
	DisplayName  string    `db:"display_name"`
}

func (u userCredentials) toEntity() entity.UserCredentials {
	return entity.UserCredentials{
		User: entity.User{
			ID:          entity.UserID(u.ID),
			Login:       u.Login,
			DisplayName: u.DisplayName,
		},
		PasswordHash: u.PasswordHash,
	}
}
