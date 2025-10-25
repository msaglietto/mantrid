package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/msaglietto/mantrid/internal/config"
	"github.com/msaglietto/mantrid/internal/logging"
	"github.com/msaglietto/mantrid/internal/paths"
	"github.com/msaglietto/mantrid/repository/json"
	"github.com/msaglietto/mantrid/service"
	"github.com/spf13/cobra"
)

var forceRemove bool

var removeAliasCmd = &cobra.Command{
	Use:   "remove [name]",
	Short: "Remove an alias",
	Long:  `Remove an existing alias by name. Prompts for confirmation unless --force flag is used.`,
	Args:  cobra.ExactArgs(1),
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

		name := args[0]
		logger.Info("removing alias", "name", name)

		// Initialize repository with proper file path
		repo := json.NewAliasRepository(fm.GetAliasFilePath())
		svc := service.NewAliasService(repo)

		// Get alias details for confirmation prompt
		if !forceRemove {
			alias, err := svc.GetAlias(ctx, name)
			if err != nil {
				logger.Error("failed to get alias", "error", err)
				return fmt.Errorf("failed to get alias: %w", err)
			}

			// Show alias details and prompt for confirmation
			fmt.Printf("Alias: %s\n", alias.Name)
			fmt.Printf("Command: %s\n", alias.Command)
			fmt.Printf("\n")

			if !confirmDelete(name) {
				logger.Info("alias removal cancelled by user", "name", name)
				fmt.Println("Removal cancelled")
				return nil
			}
		}

		// Delete the alias
		if err := svc.DeleteAlias(ctx, name); err != nil {
			logger.Error("failed to delete alias", "error", err)
			return fmt.Errorf("failed to delete alias: %w", err)
		}

		logger.Info("alias removed successfully", "name", name)
		fmt.Printf("Alias '%s' removed successfully\n", name)
		return nil
	},
}

func confirmDelete(aliasName string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Are you sure you want to remove alias '%s'? (y/N): ", aliasName)

	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}

func init() {
	removeAliasCmd.Flags().BoolVarP(&forceRemove, "force", "f", false, "Skip confirmation prompt")
	aliasCmd.AddCommand(removeAliasCmd)
}
