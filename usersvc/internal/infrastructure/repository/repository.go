package repository

import (
	"time"

	"go.elastic.co/apm"
	"gorm.io/gorm"

	"usersvc/internal/core/domain"
	"usersvc/internal/core/port"
	"usersvc/pkg/common"
)

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) port.Repository {
	return &repository{db}
}

func (r *repository) FindByID(ctx common.ServiceContextManager, id int) (*domain.User, error) {
	span, context := apm.StartSpan(ctx.GetContext(), "repository.user.FindByID", "repository")
	defer span.End()

	ctx.SetContext(context)

	var model User

	tx := r.db.WithContext(ctx.GetContext())
	if err := tx.Where("id", id).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return model.ToEntity(), nil
}

func (r *repository) FindByPhone(ctx common.ServiceContextManager, phone string) (*domain.User, error) {
	span, context := apm.StartSpan(ctx.GetContext(), "repository.user.FindByPhone", "repository")
	defer span.End()

	ctx.SetContext(context)

	var model User

	tx := r.db.WithContext(ctx.GetContext())
	if err := tx.Where("phone", phone).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return model.ToEntity(), nil
}

func (r *repository) UpdateLastLogin(ctx common.ServiceContextManager, userID int) error {
	span, context := apm.StartSpan(ctx.GetContext(), "repository.user.UpdateLastLogin", "repository")
	defer span.End()

	ctx.SetContext(context)

	tx := r.db.WithContext(ctx.GetContext())
	result := tx.Model(&User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"last_login": time.Now(),
	})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
