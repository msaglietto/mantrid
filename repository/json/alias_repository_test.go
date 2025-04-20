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
