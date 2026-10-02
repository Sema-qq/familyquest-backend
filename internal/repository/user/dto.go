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

type userDTO struct {
	ID           uuid.UUID `db:"id"`
	Login        string    `db:"login"`
	PasswordHash string    `db:"password_hash"`
	DisplayName  string    `db:"display_name"`
}

func (u userDTO) toEntity() entity.User {
	return entity.User{
		ID:           entity.UserID(u.ID),
		Login:        u.Login,
		PasswordHash: u.PasswordHash,
		DisplayName:  u.DisplayName,
	}
}
