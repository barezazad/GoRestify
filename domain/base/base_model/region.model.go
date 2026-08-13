package base_model

import (
	"time"
)

// RegionTable is used inside the repo layer for specify the table name
const (
	RegionTable = "base_regions"
)

// Region model
type Region struct {
	ID        uint       `json:"id,omitempty"`
	Name      string     `gorm:"type:varchar(100)" json:"name,omitempty" bind:"required"`
	CreatedAt *time.Time `gorm:"->;type:timestamp;not null;default:current_timestamp;" json:"created_at,omitempty"`
	UpdatedAt *time.Time `gorm:"<-:update;type:timestamp;not null;default:current_timestamp;" json:"updated_at,omitempty"`
}
