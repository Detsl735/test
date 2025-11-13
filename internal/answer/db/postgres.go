package db

import (
	"context"
	"errors"
	"fmt"
	"testTask/internal/answer"
	"testTask/pkg/logging"

	"gorm.io/gorm"
)

type repository struct {
	db     *gorm.DB
	logger *logging.Logger
}

func NewStorage(db *gorm.DB, logger *logging.Logger) answer.Storage {
	return &repository{db: db, logger: logger}
}

func (r *repository) Create(ctx context.Context, a *answer.Answer) (*answer.Answer, error) {
	if err := r.db.WithContext(ctx).Create(a).Error; err != nil {
		r.logger.Errorf("failed to create answer: %v", err)
		return nil, fmt.Errorf("create answer: %w", err)
	}
	return a, nil
}

func (r *repository) FindOne(ctx context.Context, id uint) (*answer.Answer, error) {
	var a answer.Answer
	if err := r.db.WithContext(ctx).First(&a, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.logger.Errorf("failed to find answer id=%d: %v", id, err)
		return nil, fmt.Errorf("find answer: %w", err)
	}
	return &a, nil
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&answer.Answer{}, id).Error; err != nil {
		r.logger.Errorf("failed to delete answer id=%d: %v", id, err)
		return fmt.Errorf("delete answer: %w", err)
	}
	return nil
}
