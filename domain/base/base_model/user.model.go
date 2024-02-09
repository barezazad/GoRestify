package base_model

// UserTable is used inside the repo layer for specify the table name
const (
	UserTable = "base_users"
)

// User model
type User struct {
	ID     uint   `json:"id,omitempty"`
	RoleID uint   `gorm:"index:role_id_idx;not null" json:"role_id,omitempty" bind:"required"`
	Role   string `gorm:"-:migration;->" json:"role,omitempty" table:"base_roles.name as role"`
}
