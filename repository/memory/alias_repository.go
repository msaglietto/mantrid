package memory

import (
	"context"
	"sync"

	"github.com/msaglietto/mantrid/domain"
	"github.com/msaglietto/mantrid/repository"
)

type aliasRepository struct {
	mu      sync.RWMutex
	aliases map[string]*domain.Alias
}

// NewAliasRepository creates a new in-memory alias repository.
func NewAliasRepository() repository.AliasRepository {
	return &aliasRepository{
		aliases: make(map[string]*domain.Alias),
	}
}

func (r *aliasRepository) Create(ctx context.Context, alias *domain.Alias) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.aliases[alias.Name]; exists {
		return domain.ErrAliasExists
	}

	r.aliases[alias.Name] = alias
	return nil
}

func (r *aliasRepository) FindByName(ctx context.Context, name string) (*domain.Alias, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	alias, ok := r.aliases[name]
	if !ok {
		return nil, domain.ErrAliasNotFound
	}
	return alias, nil
}

func (r *aliasRepository) List(ctx context.Context) ([]*domain.Alias, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.Alias, 0, len(r.aliases))
	for _, alias := range r.aliases {
		result = append(result, alias)
	}
	return result, nil
}

func (r *aliasRepository) Update(ctx context.Context, alias *domain.Alias) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.aliases[alias.Name]; !ok {
		return domain.ErrAliasNotFound
	}

	r.aliases[alias.Name] = alias
	return nil
}

func (r *aliasRepository) Delete(ctx context.Context, name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.aliases[name]; !ok {
		return domain.ErrAliasNotFound
	}

	delete(r.aliases, name)
	return nil
}
