package domain

import "time"

type UserMeta struct {
	At *time.Time `json:"at"`
	By *int       `json:"by"`
}
