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
