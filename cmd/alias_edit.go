package cmd

import (
	"fmt"

	"github.com/msaglietto/mantrid/internal/logging"
	"github.com/spf13/cobra"
)

var editAliasCmd = &cobra.Command{
	Use:   "edit [name] [new-command]",
	Short: "Edit an existing alias",
	Long:  `Update the command of an existing alias by name.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		application, err := appFactory(cmd.Context(), GetConfigFile())
		if err != nil {
			return err
		}

		ctx := logging.WithLogger(cmd.Context(), application.Logger)
		name, newCommand := args[0], args[1]

		application.Logger.Info("editing alias", "name", name)

		if err := application.AliasService.UpdateAlias(ctx, name, newCommand); err != nil {
			application.Logger.Error("failed to update alias", "error", err)
			return fmt.Errorf("failed to update alias: %w", err)
		}

		application.Logger.Info("alias updated successfully", "name", name)
		fmt.Printf("Alias '%s' updated successfully\n", name)
		return nil
	},
}

func init() {
	aliasCmd.AddCommand(editAliasCmd)
}
