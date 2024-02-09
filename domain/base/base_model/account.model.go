package base_model

import (
	"GoRestify/pkg/pkg_types"
	"time"
)

// AccountTable is used inside the repo layer for specify the table name
const (
	AccountTable = "base_accounts"
)

// Account model
type Account struct {
	ID        uint           `json:"id,omitempty"`
	FullName  string         `gorm:"type:varchar(45)" json:"full_name,omitempty" bind:"required"`
	Username  string         `gorm:"type:varchar(30)" json:"username,omitempty" bind:"required"`
	Password  string         `gorm:"type:varchar(200)" json:"password,omitempty" bind:"if_exist,password"`
	Email     string         `gorm:"type:varchar(60)" json:"email,omitempty" bind:"email"`
	Phone     string         `gorm:"type:varchar(20)" json:"phone,omitempty" bind:"phone"`
	Status    pkg_types.Enum `gorm:"type:varchar(25);" json:"status,omitempty" bind:"one_of=account_status"`
	Type      pkg_types.Enum `gorm:"type:varchar(25);" json:"type,omitempty" bind:"one_of=account_type"`
	CreatedAt *time.Time     `gorm:"->;type:timestamp;not null;default:current_timestamp;" json:"created_at,omitempty"`
	UpdatedAt *time.Time     `gorm:"<-:update;type:timestamp;not null;default:current_timestamp;" json:"updated_at,omitempty"`
	Token     string         `gorm:"-" json:"token"`
	Resources []string       `gorm:"-" json:"resources"`
	RoleID    uint           `gorm:"-" json:"role_id,omitempty"`
	User      User           `gorm:"-" json:"user,omitempty"`
}
