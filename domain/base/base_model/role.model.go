package base_model

import (
	"time"
)

// RoleTable is used inside the repo layer for specify the table name
const (
	RoleTable = "base_roles"
)

// Role model
type Role struct {
	ID        uint       `json:"id,omitempty"`
	Name      string     `gorm:"type:varchar(100);unique" json:"name,omitempty" bind:"required"`
	Resources string     `gorm:"type:text" json:"resources,omitempty" bind:"required"`
	CreatedAt *time.Time `gorm:"->;type:timestamp;not null;default:current_timestamp;" json:"created_at,omitempty"`
	UpdatedAt *time.Time `gorm:"<-:update;type:timestamp;not null;default:current_timestamp;" json:"updated_at,omitempty"`
}
