package json_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/msaglietto/mantrid/domain"
	"github.com/msaglietto/mantrid/repository/json"
	"github.com/stretchr/testify/assert"
)

func TestAliasRepository_List(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "aliases.json")
	repo := json.NewAliasRepository(filePath)
	ctx := context.Background()

	t.Run("list empty repository", func(t *testing.T) {
		aliases, err := repo.List(ctx)
		assert.NoError(t, err)
		assert.Empty(t, aliases)
	})

	t.Run("list with aliases", func(t *testing.T) {
		// Add some test aliases
		alias1, _ := domain.NewAlias("test1", "echo test1")
		alias2, _ := domain.NewAlias("test2", "echo test2")

		err := repo.Save(ctx, alias1)
		assert.NoError(t, err)
		err = repo.Save(ctx, alias2)
		assert.NoError(t, err)

		// List aliases
		aliases, err := repo.List(ctx)
		assert.NoError(t, err)
		assert.Len(t, aliases, 2)

		// Verify aliases content
		assert.Contains(t, []string{aliases[0].Name, aliases[1].Name}, "test1")
		assert.Contains(t, []string{aliases[0].Name, aliases[1].Name}, "test2")
	})
}

func TestAliasRepository_Update(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "aliases.json")
	repo := json.NewAliasRepository(filePath)
	ctx := context.Background()

	t.Run("update existing alias", func(t *testing.T) {
		// Create and save an alias
		alias, _ := domain.NewAlias("test", "echo original")
		err := repo.Save(ctx, alias)
		assert.NoError(t, err)

		// Update the alias
		err = alias.UpdateCommand("echo updated")
		assert.NoError(t, err)

		err = repo.Update(ctx, alias)
		assert.NoError(t, err)

		// Verify the update
		updated, err := repo.FindByName(ctx, "test")
		assert.NoError(t, err)
		assert.NotNil(t, updated)
		assert.Equal(t, "echo updated", updated.Command)
		// Compare timestamps without monotonic clock reading
		assert.True(t, alias.UpdatedAt.Truncate(0).Equal(updated.UpdatedAt.Truncate(0)))
	})

	t.Run("update non-existent alias", func(t *testing.T) {
		// Try to update an alias that doesn't exist
		alias, _ := domain.NewAlias("nonexistent", "echo test")
		err := repo.Update(ctx, alias)
		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrAliasNotFound)
	})

	t.Run("update preserves other aliases", func(t *testing.T) {
		// Create a new temp directory for this test
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "aliases.json")
		repo := json.NewAliasRepository(filePath)

		// Create multiple aliases
		alias1, _ := domain.NewAlias("alias1", "echo one")
		alias2, _ := domain.NewAlias("alias2", "echo two")
		alias3, _ := domain.NewAlias("alias3", "echo three")

		err := repo.Save(ctx, alias1)
		assert.NoError(t, err)
		err = repo.Save(ctx, alias2)
		assert.NoError(t, err)
		err = repo.Save(ctx, alias3)
		assert.NoError(t, err)

		// Update only alias2
		err = alias2.UpdateCommand("echo updated two")
		assert.NoError(t, err)
		err = repo.Update(ctx, alias2)
		assert.NoError(t, err)

		// Verify all aliases
		aliases, err := repo.List(ctx)
		assert.NoError(t, err)
		assert.Len(t, aliases, 3)

		// Find and verify each alias
		for _, a := range aliases {
			if a.Name == "alias1" {
				assert.Equal(t, "echo one", a.Command)
			} else if a.Name == "alias2" {
				assert.Equal(t, "echo updated two", a.Command)
			} else if a.Name == "alias3" {
				assert.Equal(t, "echo three", a.Command)
			}
		}
	})
}

func TestAliasRepository_Delete(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "aliases.json")
	repo := json.NewAliasRepository(filePath)
	ctx := context.Background()

	t.Run("delete existing alias", func(t *testing.T) {
		// Create and save an alias
		alias, _ := domain.NewAlias("test", "echo test")
		err := repo.Save(ctx, alias)
		assert.NoError(t, err)

		// Delete the alias
		err = repo.Delete(ctx, "test")
		assert.NoError(t, err)

		// Verify the alias is gone
		_, err = repo.FindByName(ctx, "test")
		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrAliasNotFound)
	})

	t.Run("delete non-existent alias", func(t *testing.T) {
		// Try to delete an alias that doesn't exist
		err := repo.Delete(ctx, "nonexistent")
		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrAliasNotFound)
	})

	t.Run("delete preserves other aliases", func(t *testing.T) {
		// Create a new temp directory for this test
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "aliases.json")
		repo := json.NewAliasRepository(filePath)

		// Create multiple aliases
		alias1, _ := domain.NewAlias("alias1", "echo one")
		alias2, _ := domain.NewAlias("alias2", "echo two")
		alias3, _ := domain.NewAlias("alias3", "echo three")

		err := repo.Save(ctx, alias1)
		assert.NoError(t, err)
		err = repo.Save(ctx, alias2)
		assert.NoError(t, err)
		err = repo.Save(ctx, alias3)
		assert.NoError(t, err)

		// Delete only alias2
		err = repo.Delete(ctx, "alias2")
		assert.NoError(t, err)

		// Verify remaining aliases
		aliases, err := repo.List(ctx)
		assert.NoError(t, err)
		assert.Len(t, aliases, 2)

		// Verify alias2 is gone
		_, err = repo.FindByName(ctx, "alias2")
		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrAliasNotFound)

		// Verify alias1 and alias3 still exist
		found1, err := repo.FindByName(ctx, "alias1")
		assert.NoError(t, err)
		assert.NotNil(t, found1)
		assert.Equal(t, "echo one", found1.Command)

		found3, err := repo.FindByName(ctx, "alias3")
		assert.NoError(t, err)
		assert.NotNil(t, found3)
		assert.Equal(t, "echo three", found3.Command)
	})

	t.Run("delete with empty aliases file", func(t *testing.T) {
		// Create a new temp directory for this test
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "aliases.json")
		repo := json.NewAliasRepository(filePath)

		// Try to delete from empty repository
		err := repo.Delete(ctx, "test")
		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrAliasNotFound)
	})
}
