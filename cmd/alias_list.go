// cmd/alias_list.go
package cmd

import (
	stdjson "encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/msaglietto/mantrid/internal/logging"
	"github.com/spf13/cobra"
)

var listAliasCmd = &cobra.Command{
	Use:   "list",
	Short: "List all aliases",
	Long:  `Display a list of all configured aliases with their commands and creation dates.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		application, err := appFactory(cmd.Context(), GetConfigFile())
		if err != nil {
			return err
		}

		ctx := logging.WithLogger(cmd.Context(), application.Logger)
		application.Logger.Info("listing aliases")

		// Get all aliases
		aliases, err := application.AliasService.ListAliases(ctx)
		if err != nil {
			application.Logger.Error("failed to list aliases", "error", err)
			return fmt.Errorf("failed to list aliases: %w", err)
		}

		if len(aliases) == 0 {
			// Check if JSON output is requested
			jsonOutput, _ := cmd.Flags().GetBool("json")
			if jsonOutput {
				fmt.Println("[]")
			} else {
				fmt.Println("No aliases found")
			}
			return nil
		}

		// Check if JSON output is requested
		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			output, err := stdjson.MarshalIndent(aliases, "", "  ")
			if err != nil {
				application.Logger.Error("failed to marshal aliases to JSON", "error", err)
				return fmt.Errorf("failed to marshal aliases to JSON: %w", err)
			}
			fmt.Println(string(output))
			return nil
		}

		// Initialize tabwriter for formatted output
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tCOMMAND\tCREATED\t")
		fmt.Fprintln(w, "----\t-------\t-------\t")

		for _, alias := range aliases {
			fmt.Fprintf(w, "%s\t%s\t%s\t\n",
				alias.Name,
				alias.Command,
				formatTime(alias.CreatedAt),
			)
		}

		return w.Flush()
	},
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func init() {
	aliasCmd.AddCommand(listAliasCmd)
	listAliasCmd.Flags().Bool("json", false, "Output aliases in JSON format")
}
