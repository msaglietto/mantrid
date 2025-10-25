package json

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/msaglietto/mantrid/domain"
	"github.com/msaglietto/mantrid/repository"
)

type aliasRepository struct {
	filePath string
	mu       sync.RWMutex
}

func NewAliasRepository(filePath string) repository.AliasRepository {
	return &aliasRepository{
		filePath: filePath,
	}
}

func (r *aliasRepository) Save(ctx context.Context, alias *domain.Alias) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	aliases, err := r.readAliases()
	if err != nil {
		return err
	}

	// Check for existing alias
	for i, a := range aliases {
		if a.Name == alias.Name {
			aliases[i] = alias
			return r.writeAliases(aliases)
		}
	}

	aliases = append(aliases, alias)
	return r.writeAliases(aliases)
}

func (r *aliasRepository) FindByName(ctx context.Context, name string) (*domain.Alias, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	aliases, err := r.readAliases()
	if err != nil {
		return nil, err
	}

	for _, alias := range aliases {
		if alias.Name == name {
			return alias, nil
		}
	}

	return nil, domain.ErrAliasNotFound
}

func (r *aliasRepository) readAliases() ([]*domain.Alias, error) {
	if _, err := os.Stat(r.filePath); os.IsNotExist(err) {
		return []*domain.Alias{}, nil
	}

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil, err
	}

	var aliases []*domain.Alias
	if err := json.Unmarshal(data, &aliases); err != nil {
		return nil, err
	}

	return aliases, nil
}

func (r *aliasRepository) writeAliases(aliases []*domain.Alias) error {
	data, err := json.MarshalIndent(aliases, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(r.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(r.filePath, data, 0644)
}

func (r *aliasRepository) List(ctx context.Context) ([]*domain.Alias, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	aliases, err := r.readAliases()
	if err != nil {
		return nil, fmt.Errorf("failed to read aliases: %w", err)
	}

	return aliases, nil
}

func (r *aliasRepository) Update(ctx context.Context, alias *domain.Alias) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	aliases, err := r.readAliases()
	if err != nil {
		return fmt.Errorf("failed to read aliases: %w", err)
	}

	// Find and update the alias
	found := false
	for i, a := range aliases {
		if a.Name == alias.Name {
			aliases[i] = alias
			found = true
			break
		}
	}

	if !found {
		return domain.ErrAliasNotFound
	}

	if err := r.writeAliases(aliases); err != nil {
		return fmt.Errorf("failed to write aliases: %w", err)
	}

	return nil
}

func (r *aliasRepository) Delete(ctx context.Context, name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	aliases, err := r.readAliases()
	if err != nil {
		return fmt.Errorf("failed to read aliases: %w", err)
	}

	// Find and remove the alias
	found := false
	newAliases := make([]*domain.Alias, 0, len(aliases))
	for _, a := range aliases {
		if a.Name == name {
			found = true
			// Skip this alias (effectively removing it)
			continue
		}
		newAliases = append(newAliases, a)
	}

	if !found {
		return domain.ErrAliasNotFound
	}

	if err := r.writeAliases(newAliases); err != nil {
		return fmt.Errorf("failed to write aliases: %w", err)
	}

	return nil
}
