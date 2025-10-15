package main

import (
	"fmt"
	"os"

	"crush-session-explorer/internal/cli"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "crush-md",
		Short: "Crush Session Explorer",
		Long:  "A CLI tool for exporting Crush chat sessions from SQLite databases to Markdown format",
	}

	// Add export command
	rootCmd.AddCommand(cli.ExportCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}