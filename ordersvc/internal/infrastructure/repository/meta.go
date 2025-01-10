package repository

import "time"

type CreateMeta struct {
	CreatedAt time.Time `gorm:"column:created_at;default:current_timestamp"`
	CreatedBy int       `gorm:"column:created_by;type:int"`
}

type UpdateMeta struct {
	UpdatedAt *time.Time `gorm:"column:updated_at"`
	UpdatedBy *int       `gorm:"column:updated_by;type:int"`
}

type DeleteMeta struct {
	DeletedAt *time.Time `gorm:"column:deleted_at"`
	DeletedBy *int       `gorm:"column:deleted_by;type:int"`
}

type SoftDelete struct {
	UpdateMeta
	DeleteMeta
}
