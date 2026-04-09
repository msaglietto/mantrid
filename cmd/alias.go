package cmd

import (
	"fmt"

	"github.com/msaglietto/mantrid/internal/logging"
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
		application, err := appFactory(cmd.Context(), GetConfigFile())
		if err != nil {
			return err
		}

		ctx := logging.WithLogger(cmd.Context(), application.Logger)
		name, command := args[0], args[1]

		application.Logger.Info("adding new alias", "name", name)

		if err := application.AliasService.CreateAlias(ctx, name, command); err != nil {
			application.Logger.Error("failed to create alias", "error", err)
			return fmt.Errorf("failed to create alias: %w", err)
		}

		application.Logger.Info("alias created successfully", "name", name)
		fmt.Fprintf(cmd.OutOrStdout(), "Alias '%s' created successfully\n", name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(aliasCmd)
	aliasCmd.AddCommand(addAliasCmd)
}
