package json_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/msaglietto/mantrid/domain"
	"github.com/msaglietto/mantrid/repository/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAliasRepositoryIntegration(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "aliases.json")
	repo := json.NewAliasRepository(filePath)
	ctx := context.Background()

	t.Run("save and retrieve alias", func(t *testing.T) {
		alias, err := domain.NewAlias("test", "echo test")
		require.NoError(t, err)

		err = repo.Save(ctx, alias)
		require.NoError(t, err)

		retrieved, err := repo.FindByName(ctx, "test")
		require.NoError(t, err)
		assert.Equal(t, alias.Name, retrieved.Name)
		assert.Equal(t, alias.Command, retrieved.Command)
	})

	t.Run("file persistence", func(t *testing.T) {
		// Verify file exists
		_, err := os.Stat(filePath)
		assert.NoError(t, err)

		// Verify file contents
		data, err := os.ReadFile(filePath)
		assert.NoError(t, err)
		assert.Contains(t, string(data), "test")
		assert.Contains(t, string(data), "echo test")
	})
}
