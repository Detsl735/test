package question

import (
	"context"
	"errors"
	"strings"

	"testTask/pkg/logging"
)

var (
	ErrEmptyText = errors.New("question text is empty")
	ErrNotFound  = errors.New("question not found")
)

type Service interface {
	Create(ctx context.Context, req *CreateQuestionRequest) (*Question, error)
	GetByID(ctx context.Context, id uint) (*Question, error)
	GetAll(ctx context.Context) ([]Question, error)
	Delete(ctx context.Context, id uint) error
}

type service struct {
	storage Storage
	logger  *logging.Logger
}

func NewService(storage Storage, logger *logging.Logger) Service {
	return &service{
		storage: storage,
		logger:  logger,
	}
}

func (s *service) Create(ctx context.Context, req *CreateQuestionRequest) (*Question, error) {
	text := strings.TrimSpace(req.Text)
	if text == "" {
		return nil, ErrEmptyText
	}

	q := &Question{
		Text: text,
	}

	created, err := s.storage.Create(ctx, q)
	if err != nil {
		s.logger.Errorf("failed to create question: %v", err)
		return nil, err
	}

	return created, nil
}

func (s *service) GetByID(ctx context.Context, id uint) (*Question, error) {
	q, err := s.storage.FindOne(ctx, id)
	if err != nil {
		s.logger.Errorf("failed to get question id=%d: %v", id, err)
		return nil, err
	}
	if q == nil {
		return nil, ErrNotFound
	}
	return q, nil
}

func (s *service) GetAll(ctx context.Context) ([]Question, error) {
	list, err := s.storage.FindAll(ctx)
	if err != nil {
		s.logger.Errorf("failed to list questions: %v", err)
		return nil, err
	}
	return list, nil
}

func (s *service) Delete(ctx context.Context, id uint) error {
	if err := s.storage.Delete(ctx, id); err != nil {
		s.logger.Errorf("failed to delete question id=%d: %v", id, err)
		return err
	}
	return nil
}
