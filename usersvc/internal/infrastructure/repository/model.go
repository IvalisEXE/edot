package repository

import (
	"time"

	"usersvc/internal/core/domain"
)

type User struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement"`
	ShopID    int       `gorm:"column:shop_id; null" json:"shop_id"`
	Name      string    `gorm:"column:name;size:255"`
	Phone     string    `gorm:"column:phone;size:20;unique"`
	Password  string    `gorm:"column:password;size:255"`
	Platform  string    `gorm:"column:platform"`
	LastLogin time.Time `gorm:"column:last_login;default:current_timestamp"`
	CreatedAt time.Time `gorm:"column:created_at;default:current_timestamp"`
}

func (User) TableName() string {
	return "user"
}

func (u User) ToEntity() *domain.User {
	return &domain.User{
		ID:        u.ID,
		ShopID:    u.ShopID,
		Phone:     u.Phone,
		Name:      u.Name,
		Platform:  u.Platform,
		Password:  &u.Password,
		LastLogin: &u.LastLogin,
		CreatedAt: &u.CreatedAt,
	}
}
