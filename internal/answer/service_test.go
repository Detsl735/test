package answer

import (
	"context"
	"errors"
	"testing"

	"testTask/pkg/logging"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockStorage struct {
	mock.Mock
}

func (m *mockStorage) Create(ctx context.Context, a *Answer) (*Answer, error) {
	args := m.Called(ctx, a)
	if v := args.Get(0); v != nil {
		return v.(*Answer), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockStorage) FindOne(ctx context.Context, id uint) (*Answer, error) {
	args := m.Called(ctx, id)
	if v := args.Get(0); v != nil {
		return v.(*Answer), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockStorage) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func newTestService(t *testing.T) (*service, *mockStorage) {
	t.Helper()

	logger := logging.GetLogger()
	storage := &mockStorage{}

	svc := &service{
		storage: storage,
		logger:  logger,
	}

	return svc, storage
}

func TestService_Create_InvalidQuestionID(t *testing.T) {
	svc, storage := newTestService(t)
	ctx := context.Background()

	req := &CreateAnswerRequest{
		QuestionID: 0,
		UserID:     "jh24h5",
		Text:       "test",
	}

	a, err := svc.Create(ctx, req)

	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidQuestion))
	assert.Nil(t, a)

	storage.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestService_Create_EmptyUserID(t *testing.T) {
	svc, storage := newTestService(t)
	ctx := context.Background()

	req := &CreateAnswerRequest{
		QuestionID: 1,
		UserID:     "  ",
		Text:       "test",
	}

	a, err := svc.Create(ctx, req)

	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrEmptyUserID))
	assert.Nil(t, a)

	storage.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestService_Create_EmptyText(t *testing.T) {
	svc, storage := newTestService(t)
	ctx := context.Background()

	req := &CreateAnswerRequest{
		QuestionID: 1,
		UserID:     "jh24h5",
		Text:       "   ",
	}

	a, err := svc.Create(ctx, req)

	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrEmptyText))
	assert.Nil(t, a)

	storage.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestService_Create_OK(t *testing.T) {
	svc, storage := newTestService(t)
	ctx := context.Background()

	req := &CreateAnswerRequest{
		QuestionID: 10,
		UserID:     "  jh24h5  ",
		Text:       "  test text  ",
	}

	storage.
		On("Create", mock.Anything, mock.MatchedBy(func(a *Answer) bool {
			return a.QuestionID == 10 &&
				a.UserID == "jh24h5" && // проверяем, что trim сработал
				a.Text == "test text"
		})).
		Return(&Answer{
			ID:         1,
			QuestionID: 10,
			UserID:     "jh24h5",
			Text:       "test text",
		}, nil)

	a, err := svc.Create(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, a)
	assert.Equal(t, uint(1), a.ID)
	assert.Equal(t, uint(10), a.QuestionID)
	assert.Equal(t, "jh24h5", a.UserID)
	assert.Equal(t, "test text", a.Text)

	storage.AssertExpectations(t)
}

func TestService_GetByID_NotFound(t *testing.T) {
	svc, storage := newTestService(t)
	ctx := context.Background()

	storage.
		On("FindOne", mock.Anything, uint(42)).
		Return((*Answer)(nil), nil)

	a, err := svc.GetByID(ctx, 42)

	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotFound))
	assert.Nil(t, a)

	storage.AssertExpectations(t)
}

func TestService_GetByID_StorageError(t *testing.T) {
	svc, storage := newTestService(t)
	ctx := context.Background()

	storageErr := errors.New("db is down")

	storage.
		On("FindOne", mock.Anything, uint(42)).
		Return((*Answer)(nil), storageErr)

	a, err := svc.GetByID(ctx, 42)

	require.Error(t, err)
	assert.True(t, errors.Is(err, storageErr))
	assert.Nil(t, a)

	storage.AssertExpectations(t)
}

func TestService_GetByID_OK(t *testing.T) {
	svc, storage := newTestService(t)
	ctx := context.Background()

	expected := &Answer{
		ID:         1,
		QuestionID: 10,
		UserID:     "jh24h5",
		Text:       "test",
	}

	storage.
		On("FindOne", mock.Anything, uint(1)).
		Return(expected, nil)

	a, err := svc.GetByID(ctx, 1)

	require.NoError(t, err)
	require.NotNil(t, a)
	assert.Equal(t, expected, a)

	storage.AssertExpectations(t)
}

func TestService_Delete_OK(t *testing.T) {
	svc, storage := newTestService(t)
	ctx := context.Background()

	storage.
		On("Delete", mock.Anything, uint(5)).
		Return(nil)

	err := svc.Delete(ctx, 5)

	require.NoError(t, err)
	storage.AssertExpectations(t)
}

func TestService_Delete_Error(t *testing.T) {
	svc, storage := newTestService(t)
	ctx := context.Background()

	delErr := errors.New("cannot delete")

	storage.
		On("Delete", mock.Anything, uint(5)).
		Return(delErr)

	err := svc.Delete(ctx, 5)

	require.Error(t, err)
	assert.True(t, errors.Is(err, delErr))

	storage.AssertExpectations(t)
}
