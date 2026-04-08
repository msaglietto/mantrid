package cmd

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/msaglietto/mantrid/internal/app"
	"github.com/msaglietto/mantrid/internal/config"
	"github.com/msaglietto/mantrid/internal/logging"
	"github.com/msaglietto/mantrid/repository/memory"
	"github.com/msaglietto/mantrid/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestApp creates a test environment with an in-memory app factory.
func setupTestApp(t *testing.T) *app.App {
	t.Helper()

	logger := slog.New(slog.NewTextHandler(ioDiscard{}, &slog.HandlerOptions{Level: slog.LevelDebug}))

	cfg := &config.Config{
		StorageType: "memory",
		LogLevel:    "debug",
		LogFormat:   "text",
	}

	memRepo := memory.NewAliasRepository()
	svc := service.NewAliasService(memRepo)

	application := &app.App{
		Config:       cfg,
		Logger:       logger,
		AliasService: svc,
	}

	originalFactory := appFactory
	t.Cleanup(func() {
		appFactory = originalFactory
	})

	appFactory = func(ctx context.Context, configFilePath string) (*app.App, error) {
		return application, nil
	}

	return application
}

type ioDiscard struct{}

func (ioDiscard) Write(p []byte) (int, error) { return len(p), nil }

// runCommand executes a cobra command with the given arguments and returns
// the combined output and error.
func runCommand(args ...string) (string, error) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(args)

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := rootCmd.ExecuteContext(context.Background())

	w.Close()
	os.Stdout = oldStdout

	var cmdOutput bytes.Buffer
	cmdOutput.ReadFrom(r)
	r.Close()

	return strings.TrimSpace(cmdOutput.String() + buf.String()), err
}

func TestAddAliasCommand(t *testing.T) {
	t.Run("add alias successfully", func(t *testing.T) {
		setupTestApp(t)

		output, err := runCommand("alias", "add", "test", "echo hello")
		assert.NoError(t, err)
		assert.Contains(t, output, "Alias 'test' created successfully")
	})

	t.Run("add alias missing command", func(t *testing.T) {
		setupTestApp(t)

		_, err := runCommand("alias", "add", "test")
		assert.Error(t, err)
	})

	t.Run("add alias missing name and command", func(t *testing.T) {
		setupTestApp(t)

		_, err := runCommand("alias", "add")
		assert.Error(t, err)
	})
}

func TestAddAliasDuplicate(t *testing.T) {
	setupTestApp(t)

	_, err := runCommand("alias", "add", "test", "echo hello")
	require.NoError(t, err)

	output, err := runCommand("alias", "add", "test", "echo world")
	assert.Error(t, err)
	assert.Contains(t, output, "failed to create alias")
}

func TestListAliasesCommand(t *testing.T) {
	t.Run("list empty", func(t *testing.T) {
		setupTestApp(t)

		output, err := runCommand("alias", "list")
		assert.NoError(t, err)
		assert.Contains(t, output, "No aliases found")
	})

	t.Run("list with aliases", func(t *testing.T) {
		application := setupTestApp(t)
		ctx := context.Background()

		application.AliasService.CreateAlias(ctx, "build", "go build")
		application.AliasService.CreateAlias(ctx, "test", "go test ./...")

		output, err := runCommand("alias", "list")
		assert.NoError(t, err)
		assert.Contains(t, output, "NAME")
		assert.Contains(t, output, "build")
		assert.Contains(t, output, "go build")
		assert.Contains(t, output, "test")
		assert.Contains(t, output, "go test ./...")
	})

	t.Run("list empty with json flag", func(t *testing.T) {
		setupTestApp(t)

		output, err := runCommand("alias", "list", "--json")
		assert.NoError(t, err)
		assert.Contains(t, output, "[]")
	})

	t.Run("list with json flag", func(t *testing.T) {
		application := setupTestApp(t)
		ctx := context.Background()

		application.AliasService.CreateAlias(ctx, "build", "go build")

		output, err := runCommand("alias", "list", "--json")
		assert.NoError(t, err)
		assert.Contains(t, output, "build")
		assert.Contains(t, output, "go build")
	})
}

func TestEditAliasCommand(t *testing.T) {
	t.Run("edit existing alias", func(t *testing.T) {
		application := setupTestApp(t)
		ctx := context.Background()
		application.AliasService.CreateAlias(ctx, "test", "echo original")

		output, err := runCommand("alias", "edit", "test", "echo updated")
		assert.NoError(t, err)
		assert.Contains(t, output, "Alias 'test' updated successfully")
	})

	t.Run("edit non-existent alias", func(t *testing.T) {
		setupTestApp(t)

		output, err := runCommand("alias", "edit", "nonexistent", "echo test")
		assert.Error(t, err)
		assert.Contains(t, output, "failed to update alias")
	})

	t.Run("edit missing args", func(t *testing.T) {
		setupTestApp(t)

		_, err := runCommand("alias", "edit", "test")
		assert.Error(t, err)
	})
}

func TestRemoveAliasCommand(t *testing.T) {
	t.Run("remove existing alias with force", func(t *testing.T) {
		application := setupTestApp(t)
		ctx := context.Background()
		application.AliasService.CreateAlias(ctx, "test", "echo test")

		output, err := runCommand("alias", "remove", "test", "--force")
		assert.NoError(t, err)
		assert.Contains(t, output, "Alias 'test' removed successfully")
	})

	t.Run("remove non-existent alias with force", func(t *testing.T) {
		setupTestApp(t)

		_, err := runCommand("alias", "remove", "nonexistent", "--force")
		assert.Error(t, err)
	})

	t.Run("remove missing args", func(t *testing.T) {
		setupTestApp(t)

		_, err := runCommand("alias", "remove")
		assert.Error(t, err)
	})
}

func TestDoCommand(t *testing.T) {
	t.Run("alias not found", func(t *testing.T) {
		setupTestApp(t)

		output, err := runCommand("do", "nonexistent")
		assert.Error(t, err)
		assert.Contains(t, output, "not found")
	})

	t.Run("do missing args", func(t *testing.T) {
		setupTestApp(t)

		_, err := runCommand("do")
		assert.Error(t, err)
	})
}

func TestAppFactoryError(t *testing.T) {
	originalFactory := appFactory
	t.Cleanup(func() {
		appFactory = originalFactory
	})

	appFactory = func(ctx context.Context, configFilePath string) (*app.App, error) {
		return nil, fmt.Errorf("config error")
	}

	_, err := runCommand("alias", "add", "test", "echo test")
	assert.Error(t, err)
}

func TestLoggerConfigDriven(t *testing.T) {
	tests := []struct {
		name      string
		logLevel  string
		logFormat string
	}{
		{name: "debug json", logLevel: "debug", logFormat: "json"},
		{name: "info text", logLevel: "info", logFormat: "text"},
		{name: "warn json", logLevel: "warn", logFormat: "json"},
		{name: "error text", logLevel: "error", logFormat: "text"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				LogLevel:  tt.logLevel,
				LogFormat: tt.logFormat,
			}

			logger := logging.InitLogger(cfg)
			assert.NotNil(t, logger)
		})
	}
}
