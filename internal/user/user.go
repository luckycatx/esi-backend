package user

import (
	"esi/internal/pkg/db/sqlc"

	"github.com/google/uuid"
)

// Using generated User struct from sql
type User = sqlc.User
type UUID = uuid.UUID

type Profile struct {
	Username string `json:"name"`
	Email    string `json:"email"`
}
