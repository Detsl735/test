package question

import (
	"context"
	"errors"
	"testing"

	"testTask/pkg/logging"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Create(ctx context.Context, q *Question) (*Question, error) {
	args := m.Called(ctx, q)
	if v := args.Get(0); v != nil {
		return v.(*Question), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockStorage) FindOne(ctx context.Context, id uint) (*Question, error) {
	args := m.Called(ctx, id)
	if v := args.Get(0); v != nil {
		return v.(*Question), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockStorage) FindAll(ctx context.Context) ([]Question, error) {
	args := m.Called(ctx)
	if v := args.Get(0); v != nil {
		return v.([]Question), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockStorage) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func newTestService(t *testing.T) (*service, *MockStorage) {
	t.Helper()

	logger := logging.GetLogger()
	storage := &MockStorage{}

	s := &service{
		storage: storage,
		logger:  logger,
	}

	return s, storage
}

func TestService_Create_EmptyText(t *testing.T) {
	svc, storage := newTestService(t)
	ctx := context.Background()

	storage.
		On("Create", mock.Anything, mock.Anything).
		Return(nil, nil).
		Maybe()

	q, err := svc.Create(ctx, &CreateQuestionRequest{Text: "   "})

	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrEmptyText))
	assert.Nil(t, q)

	storage.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestService_Create_OK(t *testing.T) {
	svc, storage := newTestService(t)
	ctx := context.Background()

	req := &CreateQuestionRequest{Text: "  hello  "}

	storage.
		On("Create", mock.Anything, mock.MatchedBy(func(q *Question) bool {
			return q.Text == "hello"
		})).
		Return(&Question{ID: 1, Text: "hello"}, nil)

	q, err := svc.Create(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, q)
	assert.Equal(t, uint(1), q.ID)
	assert.Equal(t, "hello", q.Text)

	storage.AssertExpectations(t)
}

func TestService_GetByID_NotFound(t *testing.T) {
	svc, storage := newTestService(t)
	ctx := context.Background()

	storage.
		On("FindOne", mock.Anything, uint(42)).
		Return((*Question)(nil), nil)

	q, err := svc.GetByID(ctx, 42)

	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotFound))
	assert.Nil(t, q)

	storage.AssertExpectations(t)
}

func TestService_GetByID_StorageError(t *testing.T) {
	svc, storage := newTestService(t)
	ctx := context.Background()

	storageErr := errors.New("db is down")
	storage.
		On("FindOne", mock.Anything, uint(42)).
		Return((*Question)(nil), storageErr)

	q, err := svc.GetByID(ctx, 42)

	require.Error(t, err)
	assert.True(t, errors.Is(err, storageErr))
	assert.Nil(t, q)

	storage.AssertExpectations(t)
}

func TestService_GetAll_OK(t *testing.T) {
	svc, storage := newTestService(t)
	ctx := context.Background()

	expected := []Question{
		{ID: 1, Text: "q1"},
		{ID: 2, Text: "q2"},
	}

	storage.
		On("FindAll", mock.Anything).
		Return(expected, nil)

	list, err := svc.GetAll(ctx)

	require.NoError(t, err)
	require.Len(t, list, 2)
	assert.Equal(t, expected, list)

	storage.AssertExpectations(t)
}

func TestService_Delete_OK(t *testing.T) {
	svc, storage := newTestService(t)
	ctx := context.Background()

	storage.
		On("Delete", mock.Anything, uint(10)).
		Return(nil)

	err := svc.Delete(ctx, 10)

	require.NoError(t, err)
	storage.AssertExpectations(t)
}

func TestService_Delete_Error(t *testing.T) {
	svc, storage := newTestService(t)
	ctx := context.Background()

	delErr := errors.New("cannot delete")
	storage.
		On("Delete", mock.Anything, uint(10)).
		Return(delErr)

	err := svc.Delete(ctx, 10)

	require.Error(t, err)
	assert.True(t, errors.Is(err, delErr))

	storage.AssertExpectations(t)
}
