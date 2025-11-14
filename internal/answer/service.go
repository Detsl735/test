package answer

import (
	"context"
	"errors"
	"strings"

	"testTask/pkg/logging"
)

var (
	ErrEmptyText       = errors.New("answer text is empty")
	ErrEmptyUserID     = errors.New("user id is empty")
	ErrInvalidQuestion = errors.New("question id is invalid")
	ErrNotFound        = errors.New("answer not found")
	ErrAlreadyAnswered = errors.New("user has already answered this question")
)

type Service interface {
	Create(ctx context.Context, req *CreateAnswerRequest) (*Answer, error)
	GetByID(ctx context.Context, id uint) (*Answer, error)
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

func (s *service) Create(ctx context.Context, req *CreateAnswerRequest) (*Answer, error) {
	text := strings.TrimSpace(req.Text)
	userID := strings.TrimSpace(req.UserID)

	if req.QuestionID == 0 {
		return nil, ErrInvalidQuestion
	}
	if userID == "" {
		return nil, ErrEmptyUserID
	}
	if text == "" {
		return nil, ErrEmptyText
	}

	existed, err := s.storage.FindByQuestionAndUser(ctx, req.QuestionID, userID)
	if err != nil {
		s.logger.Errorf("failed to check existing answer: %v", err)
		return nil, err
	}
	if existed != nil {
		return nil, ErrAlreadyAnswered
	}

	a := &Answer{
		QuestionID: req.QuestionID,
		UserID:     userID,
		Text:       text,
	}

	created, err := s.storage.Create(ctx, a)
	if err != nil {
		s.logger.Errorf("failed to create answer: %v", err)
		return nil, err
	}

	return created, nil
}

func (s *service) GetByID(ctx context.Context, id uint) (*Answer, error) {
	a, err := s.storage.FindOne(ctx, id)
	if err != nil {
		s.logger.Errorf("failed to get answer id=%d: %v", id, err)
		return nil, err
	}
	if a == nil {
		return nil, ErrNotFound
	}
	return a, nil
}

func (s *service) Delete(ctx context.Context, id uint) error {
	if err := s.storage.Delete(ctx, id); err != nil {
		s.logger.Errorf("failed to delete answer id=%d: %v", id, err)
		return err
	}
	return nil
}
