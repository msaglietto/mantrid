package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/msaglietto/mantrid/domain"
	"github.com/msaglietto/mantrid/internal/config"
	"github.com/msaglietto/mantrid/internal/logging"
	"github.com/msaglietto/mantrid/internal/paths"
	"github.com/msaglietto/mantrid/repository/json"
	"github.com/msaglietto/mantrid/service"
	"github.com/spf13/cobra"
)

var doCmd = &cobra.Command{
	Use:   "do [alias-name] [params...]",
	Short: "Execute an alias command",
	Long: `Execute a stored alias command with optional parameter substitution.

Parameter handling:
  - If alias contains placeholders ($1, $2, $@, $*), parameters are substituted
  - If alias has no placeholders, parameters are automatically appended to the end

Parameter substitution:
  $1, $2, $3...  - Positional parameters
  $@             - All parameters (space-separated)
  $*             - All parameters (same as $@)

Parameter passing:
  mantrid do <alias> [params...]      - Direct parameters
  mantrid do <alias> -- [params...]   - Parameters after -- separator
                                        (useful for passing flags like -l, --verbose)

Examples:
  # Alias without placeholders - parameters auto-append
  mantrid alias add ls "ls"
  mantrid do ls -- -la /tmp           # Executes: ls -la /tmp

  mantrid alias add dk "docker"
  mantrid do dk -- ps -a              # Executes: docker ps -a

  # Alias with placeholders - parameters substitute
  mantrid alias add greet "echo Hello, $1!"
  mantrid do greet -- World           # Executes: echo Hello, World!

  mantrid alias add deploy "kubectl apply -f $1 -n $2"
  mantrid do deploy -- app.yaml prod  # Executes: kubectl apply -f app.yaml -n prod

  # Alias with $@ - all parameters substitute in place
  mantrid alias add search "grep -r $@ ."
  mantrid do search -- "TODO"         # Executes: grep -r TODO .

WARNING: Aliases execute commands directly in your system shell.
Only create aliases for commands you trust. Parameter substitution
does not perform shell escaping - use with caution.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Extract alias name and parameters
		aliasName, params := parseDoArgs(args)

		// Bootstrap: logger, config, file manager, service
		logger := logging.GetLogger()
		ctx := logging.WithLogger(context.Background(), logger)

		cfg, err := config.Load()
		if err != nil {
			logger.Error("failed to load config", "error", err)
			return fmt.Errorf("failed to load config: %w", err)
		}

		fm := paths.NewFileManager(cfg)
		if err := fm.EnsureDirectories(); err != nil {
			logger.Error("failed to create directories", "error", err)
			return fmt.Errorf("failed to create directories: %w", err)
		}

		repo := json.NewAliasRepository(fm.GetAliasFilePath())
		svc := service.NewAliasService(repo)

		// Get the alias
		alias, err := svc.GetAlias(ctx, aliasName)
		if err != nil {
			if errors.Is(err, domain.ErrAliasNotFound) {
				logger.Error("alias not found", "name", aliasName)
				return fmt.Errorf("alias '%s' not found. Use 'mantrid alias list' to see available aliases", aliasName)
			}
			logger.Error("failed to get alias", "error", err)
			return fmt.Errorf("failed to get alias: %w", err)
		}

		logger.Info("found alias", "name", aliasName, "command", alias.Command)

		// Substitute parameters
		command := substituteParams(alias.Command, params)

		if len(params) > 0 {
			logger.Info("substituted parameters", "original", alias.Command, "final", command)
		}

		// Execute the command
		return executeCommand(ctx, command)
	},
}

// substituteParams replaces parameter placeholders in command with actual values
// Supports: $1, $2, $3, ... (positional), $@ (all params), $* (all params as string)
// If no placeholders are found and params are provided, appends params to command
func substituteParams(command string, params []string) string {
	result := command

	// If no params provided, return command as-is
	if len(params) == 0 {
		return result
	}

	// Check if command contains any parameter placeholders
	hasPlaceholders := false

	// Check for $1, $2, $3, etc.
	for i := 1; i <= len(params); i++ {
		if strings.Contains(result, fmt.Sprintf("$%d", i)) {
			hasPlaceholders = true
			break
		}
	}

	// Check for $@ or $*
	if strings.Contains(result, "$@") || strings.Contains(result, "$*") {
		hasPlaceholders = true
	}

	// If command has placeholders, do substitution
	if hasPlaceholders {
		// Replace positional parameters ($1, $2, etc.)
		for i, param := range params {
			placeholder := fmt.Sprintf("$%d", i+1)
			result = strings.ReplaceAll(result, placeholder, param)
		}

		// Replace $@ with all parameters (space-separated)
		if strings.Contains(result, "$@") {
			allParams := strings.Join(params, " ")
			result = strings.ReplaceAll(result, "$@", allParams)
		}

		// Replace $* with all parameters (same as $@)
		if strings.Contains(result, "$*") {
			allParams := strings.Join(params, " ")
			result = strings.ReplaceAll(result, "$*", allParams)
		}
	} else {
		// No placeholders found - auto-append parameters
		allParams := strings.Join(params, " ")
		result = result + " " + allParams
	}

	return result
}

// parseDoArgs extracts alias name and parameters from command arguments
// Handles both direct parameters and -- separator pattern
func parseDoArgs(args []string) (aliasName string, params []string) {
	if len(args) == 0 {
		return "", []string{}
	}

	aliasName = args[0]

	// Check if -- separator is used
	if len(args) > 1 && args[1] == "--" {
		// Format: mantrid do <alias-name> -- [params...]
		params = args[2:] // Everything after --
	} else {
		// Format: mantrid do <alias-name> [params...]
		params = args[1:]
	}

	return aliasName, params
}

// executeCommand runs the command in the system shell
// Returns the exit code of the executed command
func executeCommand(ctx context.Context, command string) error {
	logger := logging.FromContext(ctx)

	var cmd *exec.Cmd

	// Platform-specific shell selection
	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, "cmd", "/C", command)
	} else {
		// Unix-like systems (Linux, macOS)
		cmd = exec.CommandContext(ctx, "sh", "-c", command)
	}

	// Connect stdin/stdout/stderr to current process
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	logger.Info("executing alias command", "command", command)

	// Run the command
	if err := cmd.Run(); err != nil {
		// Check if it's an exit error (non-zero exit code)
		if exitErr, ok := err.(*exec.ExitError); ok {
			logger.Warn("command exited with error", "exit_code", exitErr.ExitCode())
			os.Exit(exitErr.ExitCode())
		}
		return fmt.Errorf("failed to execute command: %w", err)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(doCmd)
}
