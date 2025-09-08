package model

import (
	"time"
)

type Item struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Description string    `gorm:"type:varchar(255)" json:"description"`
	Category    string    `gorm:"type:varchar(100)" json:"category"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (i *Item) TableName() string {
	return "items"
}
