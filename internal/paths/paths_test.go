package paths_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/msaglietto/mantrid/internal/config"
	"github.com/msaglietto/mantrid/internal/paths"
	"github.com/stretchr/testify/assert"
)

func TestFileManager(t *testing.T) {
	t.Run("default path", func(t *testing.T) {
		fm := paths.NewFileManager(nil)
		homeDir, _ := os.UserHomeDir()
		expected := filepath.Join(homeDir, ".mantrid", "aliases.json")
		assert.Equal(t, expected, fm.GetAliasFilePath())
	})

	t.Run("configured path", func(t *testing.T) {
		cfg := &config.Config{
			AliasFile: "/custom/path/aliases.json",
		}
		fm := paths.NewFileManager(cfg)
		assert.Equal(t, "/custom/path/aliases.json", fm.GetAliasFilePath())
	})

	t.Run("ensure directories", func(t *testing.T) {
		// Create temporary directory for test
		tmpDir := t.TempDir()
		cfg := &config.Config{
			AliasFile: filepath.Join(tmpDir, "mantrid", "aliases.json"),
		}

		fm := paths.NewFileManager(cfg)
		err := fm.EnsureDirectories()
		assert.NoError(t, err)

		// Check if directory was created
		dirInfo, err := os.Stat(filepath.Join(tmpDir, "mantrid"))
		assert.NoError(t, err)
		assert.True(t, dirInfo.IsDir())
	})
}
