package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"geminictl/internal/registry"
	"geminictl/internal/scanner"
	"geminictl/internal/tui"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the status of Gemini CLI projects and sessions",
	Run: func(cmd *cobra.Command, args []string) {
		reg, err := registry.NewRegistry()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing registry: %v\n", err)
			os.Exit(1)
		}
		if err := reg.Load(); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading registry: %v\n", err)
			os.Exit(1)
		}

		scan, err := scanner.NewScanner()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing scanner: %v\n", err)
			os.Exit(1)
		}

		projects, err := scan.Scan()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error scanning sessions: %v\n", err)
			os.Exit(1)
		}

		m := tui.NewModel(projects, reg)
		p := tea.NewProgram(m, tea.WithAltScreen())

		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}