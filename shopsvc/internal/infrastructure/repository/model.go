package repository

import (
	"time"
)

type Shop struct {
	ID        int        `gorm:"column:id;primaryKey;autoIncrement"`
	Name      string     `gorm:"column:name;size:255"`
	CreatedAt *time.Time `gorm:"column:created_at;default:current_timestamp"`
}

func (Shop) TableName() string {
	return "shop"
}
