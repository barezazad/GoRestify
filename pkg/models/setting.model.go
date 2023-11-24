package models

import "GoRestify/pkg/pkg_types"

// SettingTable is used inside the repo layer for specify the table name
const (
	SettingTable = "base_settings"
)

// Setting model
type Setting struct {
	ID          uint              `json:"id,omitempty"`
	Property    pkg_types.Setting `gorm:"not null;unique" json:"property,omitempty"`
	Value       string            `gorm:"type:text" json:"value,omitempty"`
	Description string            `gorm:"type:varchar(255)" json:"description,omitempty"`
}

// ColumnsSetting .
var ColumnsSetting = []string{"id", "property", "value", "description"}
