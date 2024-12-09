package paths

import (
	"os"
	"path/filepath"

	"github.com/msaglietto/mantrid/internal/config"
)

// FileManager handles all file path operations in the application
type FileManager struct {
	config *config.Config
}

func NewFileManager(cfg *config.Config) *FileManager {
	return &FileManager{
		config: cfg,
	}
}

func (fm *FileManager) GetAliasFilePath() string {
	// If configured in config, use that path
	if fm.config != nil && fm.config.AliasFile != "" {
		return fm.config.AliasFile
	}

	// Otherwise, use default path in user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if can't get home directory
		return filepath.Join(".", ".mantrid", "aliases.json")
	}

	return filepath.Join(homeDir, ".mantrid", "aliases.json")
}

func (fm *FileManager) EnsureDirectories() error {
	dir := filepath.Dir(fm.GetAliasFilePath())
	return os.MkdirAll(dir, 0755)
}
