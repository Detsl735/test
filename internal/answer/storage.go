package answer

import "context"

type Storage interface {
	Create(ctx context.Context, a *Answer) (*Answer, error)
	FindOne(ctx context.Context, id uint) (*Answer, error)
	Delete(ctx context.Context, id uint) error
	FindByQuestionAndUser(ctx context.Context, questionID uint, userID string) (*Answer, error)
}
