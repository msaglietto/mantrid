package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/msaglietto/mantrid/internal/config"
	"github.com/msaglietto/mantrid/internal/logging"
	"github.com/msaglietto/mantrid/internal/paths"
	"github.com/msaglietto/mantrid/repository"
	jsonrepo "github.com/msaglietto/mantrid/repository/json"
	"github.com/msaglietto/mantrid/repository/memory"
	"github.com/msaglietto/mantrid/service"
)

// App holds all application dependencies, acting as a DI container.
type App struct {
	Config       *config.Config
	Logger       *slog.Logger
	FileManager  *paths.FileManager
	AliasService service.AliasService
}

// New creates a new App instance with all dependencies initialized.
// It accepts an optional config file path (from --config flag).
func New(ctx context.Context, configFilePath ...string) (*App, error) {
	// Load configuration
	cfg, err := config.Load(configFilePath...)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize logger from config
	logger := logging.InitLogger(cfg)
	ctx = logging.WithLogger(ctx, logger)

	// Initialize file manager
	fm := paths.NewFileManager(cfg)
	if err := fm.EnsureDirectories(); err != nil {
		return nil, fmt.Errorf("failed to create directories: %w", err)
	}

	// Initialize repository based on config
	repo := newRepository(cfg, fm)

	// Initialize service
	svc := service.NewAliasService(repo)

	return &App{
		Config:       cfg,
		Logger:       logger,
		FileManager:  fm,
		AliasService: svc,
	}, nil
}

// newRepository creates the appropriate repository based on config.
func newRepository(cfg *config.Config, fm *paths.FileManager) repository.AliasRepository {
	switch cfg.StorageType {
	case "memory":
		return memory.NewAliasRepository()
	default:
		return jsonrepo.NewAliasRepository(fm.GetAliasFilePath())
	}
}
