package repository

import (
	"gorm.io/gorm"

	"shopsvc/internal/core/port"
)

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) port.Repository {
	return &repository{db}
}
