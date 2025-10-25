package repository

import (
	"context"

	"github.com/msaglietto/mantrid/domain"
)

type AliasRepository interface {
	Save(ctx context.Context, alias *domain.Alias) error
	FindByName(ctx context.Context, name string) (*domain.Alias, error)
	List(ctx context.Context) ([]*domain.Alias, error)
	Update(ctx context.Context, alias *domain.Alias) error
	Delete(ctx context.Context, name string) error
}
