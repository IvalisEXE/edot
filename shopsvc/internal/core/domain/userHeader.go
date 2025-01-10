package domain

import "time"

type UserHeader struct {
	ID        int        `json:"id"`
	Email     string     `json:"email"`
	Name      string     `json:"name,omitempty"`
	Phone     string     `json:"phone"`
	Password  *string    `json:"-"`
	LastLogin *time.Time `json:"last_login,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}
