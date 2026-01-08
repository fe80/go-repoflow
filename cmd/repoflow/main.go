package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/fe80/go-repoflow/internal/cli"
	"github.com/fe80/go-repoflow/internal/factory"
	"github.com/fe80/go-repoflow/pkg/config"
)

var (
	debug  bool
	output string
	utils  factory.Utils
)

func main() {
	cfg, err := config.Load("")
	if err != nil {
		fmt.Printf("Fail to load configuration: %v\n", err)
		os.Exit(1)
	}
	rootCmd := &cobra.Command{Use: "repoflow"}
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug logging")
	rootCmd.PersistentFlags().StringVar(&output, "output", "text", "Define output (text, yaml, json)")
	rootCmd.RegisterFlagCompletionFunc("output", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"text", "yaml", "json"}, cobra.ShellCompDirectiveNoFileComp
	})

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		level := slog.LevelInfo
		if debug {
			level = slog.LevelDebug
		}

		handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
		logger := slog.New(handler)
		slog.SetDefault(logger)
		utils.Logger = logger
		utils.Cfg = cfg
		utils.Output = output

		return nil
	}

	rootCmd.AddCommand(cli.WorkspaceCmd(&utils))
	rootCmd.AddCommand(cli.RepositoryCmd(&utils))

	if err := rootCmd.Execute(); err != nil {
		slog.Debug("Error", "error", err)
		os.Exit(1)
	}
}
