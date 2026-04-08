package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/msaglietto/mantrid/internal/logging"
	"github.com/spf13/cobra"
)

var forceRemove bool

var removeAliasCmd = &cobra.Command{
	Use:   "remove [name]",
	Short: "Remove an alias",
	Long:  `Remove an existing alias by name. Prompts for confirmation unless --force flag is used.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		application, err := appFactory(cmd.Context(), GetConfigFile())
		if err != nil {
			return err
		}

		ctx := logging.WithLogger(cmd.Context(), application.Logger)
		name := args[0]

		application.Logger.Info("removing alias", "name", name)

		// Get alias details for confirmation prompt
		if !forceRemove {
			alias, err := application.AliasService.GetAlias(ctx, name)
			if err != nil {
				application.Logger.Error("failed to get alias", "error", err)
				return fmt.Errorf("failed to get alias: %w", err)
			}

			// Show alias details and prompt for confirmation
			fmt.Fprintf(cmd.OutOrStdout(), "Alias: %s\n", alias.Name)
			fmt.Fprintf(cmd.OutOrStdout(), "Command: %s\n", alias.Command)
			fmt.Fprintf(cmd.OutOrStdout(), "\n")

			if !confirmDelete(name) {
				application.Logger.Info("alias removal cancelled by user", "name", name)
				fmt.Fprintln(cmd.OutOrStdout(), "Removal cancelled")
				return nil
			}
		}

		// Delete the alias
		if err := application.AliasService.DeleteAlias(ctx, name); err != nil {
			application.Logger.Error("failed to delete alias", "error", err)
			return fmt.Errorf("failed to delete alias: %w", err)
		}

		application.Logger.Info("alias removed successfully", "name", name)
		fmt.Fprintf(cmd.OutOrStdout(), "Alias '%s' removed successfully\n", name)
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
