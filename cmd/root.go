package cmd

import (
	"context"

	"github.com/msaglietto/mantrid/internal/app"
	"github.com/spf13/cobra"
)

var cfgFile string

// appFactory is the function used to create an App instance.
// It can be overridden in tests to inject mock dependencies.
var appFactory = defaultAppFactory

func defaultAppFactory(ctx context.Context, configFilePath string) (*app.App, error) {
	if configFilePath != "" {
		return app.New(ctx, configFilePath)
	}
	return app.New(ctx)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mantrid",
	Short: "Mantrid is a CLI tool for managing aliases and dotfiles",
	Long:  `A command-line application written in Go that helps you store and run aliases, and manage your dotfile configurations.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(ctx context.Context) error {
	return rootCmd.ExecuteContext(ctx)
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mantrid/config.yaml)")
}

// GetConfigFile returns the config file path specified via --config flag
func GetConfigFile() string {
	return cfgFile
}
