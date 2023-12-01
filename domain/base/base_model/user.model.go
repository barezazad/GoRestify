package base_model

import (
	"GoRestify/pkg/pkg_types"
	"time"
)

// UserTable is used inside the repo layer for specify the table name
const (
	UserTable = "base_users"
)

// User model
type User struct {
	ID        uint           `json:"id,omitempty"`
	RoleID    uint           `gorm:"index:role_id_idx;not null" json:"role_id,omitempty" bind:"required"`
	Username  string         `gorm:"type:varchar(30)" json:"username,omitempty" bind:"required"`
	FullName  string         `gorm:"type:varchar(45)" json:"full_name,omitempty" bind:"required"`
	Password  string         `gorm:"type:varchar(200)" json:"password,omitempty" bind:"if_exist,password"`
	Email     string         `gorm:"type:varchar(60)" json:"email,omitempty" bind:"email"`
	Phone     string         `gorm:"type:varchar(20)" json:"phone,omitempty" bind:"phone"`
	Role      string         `gorm:"-:migration;->" json:"role,omitempty" table:"base_roles.name as role"`
	Token     string         `gorm:"-" json:"token"`
	Resources []string       `gorm:"-" json:"resources"`
	Status    pkg_types.Enum `gorm:"type:varchar(20)" json:"status,omitempty" bind:"required"`
	CreatedAt *time.Time     `gorm:"->;type:timestamp;not null;default:current_timestamp;" json:"created_at,omitempty"`
	UpdatedAt *time.Time     `gorm:"<-:update;type:timestamp;not null;default:current_timestamp;" json:"updated_at,omitempty"`
}
