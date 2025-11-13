package question

import "context"

type Storage interface {
	Create(ctx context.Context, q *Question) (*Question, error)
	FindOne(ctx context.Context, id uint) (*Question, error)
	FindAll(ctx context.Context) ([]Question, error)
	Delete(ctx context.Context, id uint) error
}
