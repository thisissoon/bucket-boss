package main

import (
	"io"
	"os"

	"github.com/thisissoon/bucket-boss/internal/config"
	"github.com/thisissoon/bucket-boss/internal/version"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	configkit "go.soon.build/kit/config"
)

// Default logger
var log zerolog.Logger

// Global app configuration
var cfg config.Config

// Application entry point
func main() {
	cmd := bucketbossCmd()
	if err := cmd.Execute(); err != nil {
		log.Error().Err(err).Msg("exiting from fatal error")
		os.Exit(1)
	}
}

// New constructs a new CLI interface for execution
func bucketbossCmd() *cobra.Command {
	var configPath string
	cmd := &cobra.Command{
		Use:           "bucket-boss",
		Short:         "Run the service",
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Load config
			var err error
			cfg, err = config.New(
				configkit.WithFile(configPath),
				configkit.BindFlag("log.console", cmd.Flag("console")),
				configkit.BindFlag("log.verbose", cmd.Flag("verbose")),
			)
			if err != nil {
				return err
			}
			// Setup default logger
			log = initLogger(cfg.Log)
			return nil
		},
		RunE: bucketbossRun,
	}
	// Global flags
	pflags := cmd.PersistentFlags()
	pflags.StringVarP(&configPath, "config", "c", "", "path to configuration file (default is $HOME/.config/bucket-boss.toml)")
	pflags.Bool("console", false, "use console log writer")
	pflags.BoolP("verbose", "v", false, "verbose logging")
	// Add sub commands
	cmd.AddCommand(versionCmd())
	cmd.AddCommand(purgeCmd())
	return cmd
}

// bucket-bossRun is executed when the CLI executes
// the bucket-boss command
func bucketbossRun(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}

// initLogger constructs a default logger from config
func initLogger(c config.Log) zerolog.Logger {
	// Set logger level field to severity for stack driver support
	zerolog.LevelFieldName = "severity"
	var w io.Writer = os.Stdout
	if c.Console {
		w = zerolog.ConsoleWriter{
			Out: os.Stdout,
		}
	}
	// Parse level from config
	lvl, err := zerolog.ParseLevel(c.Level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}
	// Override level with verbose
	if c.Verbose {
		lvl = zerolog.DebugLevel
	}
	return zerolog.New(w).Level(lvl).With().Fields(map[string]interface{}{
		"version": version.Version,
		"app":     config.APP_NAME,
	}).Timestamp().Logger()
}
