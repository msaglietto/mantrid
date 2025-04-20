// cmd/alias_list.go
package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/msaglietto/mantrid/internal/config"
	"github.com/msaglietto/mantrid/internal/logging"
	"github.com/msaglietto/mantrid/internal/paths"
	"github.com/msaglietto/mantrid/repository/json"
	"github.com/msaglietto/mantrid/service"
	"github.com/spf13/cobra"
)

var listAliasCmd = &cobra.Command{
	Use:   "list",
	Short: "List all aliases",
	Long:  `Display a list of all configured aliases with their commands and creation dates.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := logging.GetLogger()
		ctx := logging.WithLogger(context.Background(), logger)

		logger.Info("listing aliases")

		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			logger.Error("failed to load config", "error", err)
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Initialize file manager
		fm := paths.NewFileManager(cfg)

		// Initialize repository and service
		repo := json.NewAliasRepository(fm.GetAliasFilePath())
		svc := service.NewAliasService(repo)

		// Get all aliases
		aliases, err := svc.ListAliases(ctx)
		if err != nil {
			logger.Error("failed to list aliases", "error", err)
			return fmt.Errorf("failed to list aliases: %w", err)
		}

		if len(aliases) == 0 {
			fmt.Println("No aliases found")
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
}
