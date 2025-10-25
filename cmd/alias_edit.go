package cmd

import (
	"context"
	"fmt"

	"github.com/msaglietto/mantrid/internal/config"
	"github.com/msaglietto/mantrid/internal/logging"
	"github.com/msaglietto/mantrid/internal/paths"
	"github.com/msaglietto/mantrid/repository/json"
	"github.com/msaglietto/mantrid/service"
	"github.com/spf13/cobra"
)

var editAliasCmd = &cobra.Command{
	Use:   "edit [name] [new-command]",
	Short: "Edit an existing alias",
	Long:  `Update the command of an existing alias by name.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get logger
		logger := logging.GetLogger()
		ctx := logging.WithLogger(context.Background(), logger)

		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			logger.Error("failed to load config", "error", err)
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Initialize file manager
		fm := paths.NewFileManager(cfg)
		if err := fm.EnsureDirectories(); err != nil {
			logger.Error("failed to create directories", "error", err)
			return fmt.Errorf("failed to create directories: %w", err)
		}

		name, newCommand := args[0], args[1]
		logger.Info("editing alias", "name", name)

		// Initialize repository with proper file path
		repo := json.NewAliasRepository(fm.GetAliasFilePath())
		svc := service.NewAliasService(repo)

		if err := svc.UpdateAlias(ctx, name, newCommand); err != nil {
			logger.Error("failed to update alias", "error", err)
			return fmt.Errorf("failed to update alias: %w", err)
		}

		logger.Info("alias updated successfully", "name", name)
		fmt.Printf("Alias '%s' updated successfully\n", name)
		return nil
	},
}

func init() {
	aliasCmd.AddCommand(editAliasCmd)
}
