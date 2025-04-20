package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/msaglietto/mantrid/domain"
	"github.com/msaglietto/mantrid/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAliasRepository struct {
	mock.Mock
}

func cleanupMock(t *testing.T, m *MockAliasRepository) {
	t.Cleanup(func() {
		m.ExpectedCalls = nil
		m.Calls = nil
	})
}

func (m *MockAliasRepository) Save(ctx context.Context, alias *domain.Alias) error {
	args := m.Called(ctx, alias)
	return args.Error(0)
}

func (m *MockAliasRepository) FindByName(ctx context.Context, name string) (*domain.Alias, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Alias), args.Error(1)
}

func (m *MockAliasRepository) List(ctx context.Context) ([]*domain.Alias, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Alias), args.Error(1)
}

func TestCreateAlias(t *testing.T) {
	mockRepo := new(MockAliasRepository)
	service := service.NewAliasService(mockRepo)
	ctx := context.Background()

	t.Run("create new alias", func(t *testing.T) {
		cleanupMock(t, mockRepo)

		mockRepo.On("FindByName", ctx, "test").Return(nil, domain.ErrAliasNotFound)
		mockRepo.On("Save", ctx, mock.AnythingOfType("*domain.Alias")).Return(nil)

		err := service.CreateAlias(ctx, "test", "echo test")
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)

		// Clean up mock after this test case
		mockRepo.ExpectedCalls = nil
		mockRepo.Calls = nil
	})

	t.Run("alias already exists", func(t *testing.T) {
		cleanupMock(t, mockRepo)

		existingAlias, _ := domain.NewAlias("test", "echo test")
		mockRepo.On("FindByName", ctx, "test").Return(existingAlias, nil)

		err := service.CreateAlias(ctx, "test", "echo test")
		assert.Equal(t, domain.ErrAliasExists, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestListAliases(t *testing.T) {
	mockRepo := new(MockAliasRepository)
	service := service.NewAliasService(mockRepo)
	ctx := context.Background()

	t.Run("list aliases successfully", func(t *testing.T) {
		cleanupMock(t, mockRepo)

		expectedAliases := []*domain.Alias{
			{
				Name:      "test1",
				Command:   "echo test1",
				CreatedAt: time.Now(),
			},
			{
				Name:      "test2",
				Command:   "echo test2",
				CreatedAt: time.Now(),
			},
		}

		mockRepo.On("List", ctx).Return(expectedAliases, nil)

		aliases, err := service.ListAliases(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedAliases, aliases)
		assert.Len(t, aliases, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("list aliases with error", func(t *testing.T) {
		cleanupMock(t, mockRepo)

		expectedErr := errors.New("failed to list aliases")
		mockRepo.On("List", ctx).Return(nil, expectedErr)

		aliases, err := service.ListAliases(ctx)
		assert.Error(t, err)
		assert.Nil(t, aliases)
		mockRepo.AssertExpectations(t)
	})
}
