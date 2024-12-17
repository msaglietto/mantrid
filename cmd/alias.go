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

var aliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "Manage aliases",
}

var addAliasCmd = &cobra.Command{
	Use:   "add [name] [command]",
	Short: "Add a new alias",
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

		name, command := args[0], args[1]
		logger.Info("adding new alias", "name", name)

		// Initialize repository with proper file path
		repo := json.NewAliasRepository(fm.GetAliasFilePath())
		svc := service.NewAliasService(repo)

		if err := svc.CreateAlias(ctx, name, command); err != nil {
			logger.Error("failed to create alias", "error", err)
			return fmt.Errorf("failed to create alias: %w", err)
		}

		logger.Info("alias created successfully", "name", name)
		fmt.Printf("Alias '%s' created successfully\n", name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(aliasCmd)
	aliasCmd.AddCommand(addAliasCmd)
}
