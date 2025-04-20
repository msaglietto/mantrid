package service_test

import (
	"context"
	"testing"

	"github.com/msaglietto/mantrid/domain"
	"github.com/msaglietto/mantrid/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAliasRepository struct {
	mock.Mock
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

func TestCreateAlias(t *testing.T) {
	mockRepo := new(MockAliasRepository)
	service := service.NewAliasService(mockRepo)
	ctx := context.Background()

	t.Run("create new alias", func(t *testing.T) {
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
		existingAlias, _ := domain.NewAlias("test", "echo test")
		mockRepo.On("FindByName", ctx, "test").Return(existingAlias, nil)

		err := service.CreateAlias(ctx, "test", "echo test")
		assert.Equal(t, domain.ErrAliasExists, err)
		mockRepo.AssertExpectations(t)
		// Clean up mock after this test case
		mockRepo.ExpectedCalls = nil
		mockRepo.Calls = nil
	})
}
