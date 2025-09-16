package model

import (
	"time"
)

// UserInformation represents user information stored in JSON field
type UserInformation struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
}

type Task struct {
	ID              uint             `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          uint             `gorm:"not null" json:"user_id"`
	UserInformation *UserInformation `gorm:"type:jsonb;serializer:json" json:"user_information"`
	Client          string           `gorm:"type:varchar(255)" json:"client"`
	Job             string           `gorm:"type:varchar(255)" json:"job"`
	Item            string           `gorm:"type:varchar(255)" json:"item"`
	Role            string           `gorm:"type:varchar(100)" json:"role"`
	Note            string           `gorm:"type:text" json:"note"`
	OT              string           `gorm:"type:varchar(255)" json:"ot"`
	Volume          int              `gorm:"type:int" json:"volume"`
	StartedAt       *time.Time       `gorm:"type:timestamp" json:"started_at"`
	EndedAt         *time.Time       `gorm:"type:timestamp" json:"ended_at"`
	WorkTime        int              `gorm:"type:int" json:"work_time"`
	Minute          int              `gorm:"type:int" json:"minute"`
	Status          string           `gorm:"type:varchar(50)" json:"status"`
	Delivery        bool             `gorm:"default:false" json:"delivery"`
	CreatedAt       time.Time        `gorm:"autoCreateTime" json:"created_at"`
	CustomCreatedAt *time.Time       `gorm:"type:timestamp" json:"custom_created_at"`
	UpdatedAt       time.Time        `gorm:"autoUpdateTime" json:"updated_at"`
}

func (t *Task) TableName() string {
	return "tasks"
}
