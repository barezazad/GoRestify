package models

import (
	"gorm.io/gorm"
)

const (
	// ActivityTable is used inside the repo layer
	ActivityTable = "base_activities"
)

// Activity model
type Activity struct {
	gorm.Model
	Event      string `gorm:"index:event_idx" json:"event"`
	OperatorID uint   `gorm:"index:operator_idx" json:"operator_id"`
	Username   string `gorm:"index:username_idx" json:"username"`
	IP         string `json:"ip"`
	URI        string `gorm:"type:text" json:"uri"`
	Before     string `gorm:"type:text" json:"before"`
	After      string `gorm:"type:text" json:"after"`
}

// ColumnsActivity .
var ColumnsActivity = []string{
	"base_activities.id", "base_activities.created_at", "base_activities.deleted_at",
	"base_activities.updated_at", "base_activities.event", "base_activities.operator_id",
	"base_activities.username", "base_activities.ip", "base_activities.uri",
	"base_activities.before", "base_activities.after",
}
