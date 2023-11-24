package pkg_types

import (
	"time"

	"GoRestify/pkg/dictionary"

	"github.com/golang-jwt/jwt/v5"
)

// Enum is used for define all types
type Enum string

// GormCol is a same as model.gorm, we use our name if in future customize it doesn't face problem
type GormCol struct {
	ID        uint       `gorm:"primary_key" json:"id,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// JWTClaims for JWT
type JWTClaims struct {
	UserID   uint            `json:"user_id"`
	Username string          `json:"username"`
	Email    string          `json:"email"`
	Phone    string          `json:"phone"`
	Lang     dictionary.Lang `json:"language"`
	jwt.RegisteredClaims
}
