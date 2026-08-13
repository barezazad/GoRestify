package base_model

import (
	"time"
)

// CityTable is used inside the repo layer
const (
	CityTable = "base_cities"
)

// City model
type City struct {
	ID        uint       `json:"id,omitempty"`
	RegionID  uint       `gorm:"index:region_id_idx;not null" json:"region_id,omitempty" bind:"required"`
	Name      string     `gorm:"type:varchar(30);not null;unique" json:"name,omitempty" bind:"required"`
	CreatedAt *time.Time `gorm:"->;type:timestamp;not null;default:current_timestamp;" json:"created_at,omitempty"`
	UpdatedAt *time.Time `gorm:"<-:update;type:timestamp;not null;default:current_timestamp;" json:"updated_at,omitempty"`
}
