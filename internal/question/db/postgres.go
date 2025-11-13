package db

import (
	"context"
	"errors"
	"fmt"
	"testTask/internal/question"
	"testTask/pkg/logging"

	"gorm.io/gorm"
)

type repository struct {
	db     *gorm.DB
	logger *logging.Logger
}

func NewStorage(db *gorm.DB, logger *logging.Logger) question.Storage {
	return &repository{db: db, logger: logger}
}

func (r *repository) Create(ctx context.Context, q *question.Question) (*question.Question, error) {
	if err := r.db.WithContext(ctx).Create(q).Error; err != nil {
		r.logger.Errorf("failed to create question: %v", err)
		return nil, fmt.Errorf("create question: %w", err)
	}
	return q, nil
}

func (r *repository) FindOne(ctx context.Context, id uint) (*question.Question, error) {
	var q question.Question

	if err := r.db.WithContext(ctx).First(&q, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.logger.Errorf("failed to find question id=%d: %v", id, err)
		return nil, fmt.Errorf("find question: %w", err)
	}

	return &q, nil
}

func (r *repository) FindAll(ctx context.Context) ([]question.Question, error) {
	var list []question.Question

	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Find(&list).Error; err != nil {
		r.logger.Errorf("failed to list questions: %v", err)
		return nil, fmt.Errorf("list questions: %w", err)
	}

	return list, nil
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).
		Delete(&question.Question{}, id).Error; err != nil {
		r.logger.Errorf("failed to delete question id=%d: %v", id, err)
		return fmt.Errorf("delete question: %w", err)
	}
	return nil
}
