package service

import (
	"context"

	"github.com/msaglietto/mantrid/domain"
	"github.com/msaglietto/mantrid/repository"
)

type AliasService interface {
	CreateAlias(ctx context.Context, name, command string) error
	GetAlias(ctx context.Context, name string) (*domain.Alias, error)
	ListAliases(ctx context.Context) ([]*domain.Alias, error)
	UpdateAlias(ctx context.Context, name, newCommand string) error
	DeleteAlias(ctx context.Context, name string) error
}

type aliasService struct {
	repo repository.AliasRepository
}

func NewAliasService(repo repository.AliasRepository) AliasService {
	return &aliasService{
		repo: repo,
	}
}

func (s *aliasService) CreateAlias(ctx context.Context, name, command string) error {
	alias, err := domain.NewAlias(name, command)
	if err != nil {
		return err
	}

	if _, err := s.repo.FindByName(ctx, name); err == nil {
		return domain.ErrAliasExists
	}

	return s.repo.Save(ctx, alias)
}

func (s *aliasService) GetAlias(ctx context.Context, name string) (*domain.Alias, error) {
	return s.repo.FindByName(ctx, name)
}

func (s *aliasService) ListAliases(ctx context.Context) ([]*domain.Alias, error) {
	return s.repo.List(ctx)
}

func (s *aliasService) UpdateAlias(ctx context.Context, name, newCommand string) error {
	// Validate inputs
	if name == "" {
		return domain.ErrEmptyAliasName
	}
	if newCommand == "" {
		return domain.ErrEmptyAliasCommand
	}

	// Find the existing alias
	alias, err := s.repo.FindByName(ctx, name)
	if err != nil {
		return err
	}

	// Update the command
	if err := alias.UpdateCommand(newCommand); err != nil {
		return err
	}

	// Save the updated alias
	return s.repo.Update(ctx, alias)
}

func (s *aliasService) DeleteAlias(ctx context.Context, name string) error {
	// Validate input
	if name == "" {
		return domain.ErrEmptyAliasName
	}

	// Delete the alias
	return s.repo.Delete(ctx, name)
}
