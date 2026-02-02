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

var resetRegistry bool

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the status of Gemini CLI projects and sessions",
	Run: func(cmd *cobra.Command, args []string) {
		reg, err := registry.NewRegistry()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing registry: %v\n", err)
			os.Exit(1)
		}

		if resetRegistry {
			reg.Clear()
			_ = reg.Save()
		} else {
			if err := reg.Load(); err != nil {
				fmt.Fprintf(os.Stderr, "Error loading registry: %v\n", err)
				os.Exit(1)
			}
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

		// Auto-register current directory if it matches a scanned project
		cwd, err := os.Getwd()
		if err == nil {
			currentID, err := registry.CalculateProjectID(cwd)
			if err == nil {
				for _, p := range projects {
					if p.ID == currentID {
						reg.AddProject(currentID, cwd)
						_ = reg.Save() // Ignore error on auto-save
						break
					}
				}
			}
		}

		m := tui.NewModel(projects, reg, scan)
		p := tea.NewProgram(m, tea.WithAltScreen())

		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	statusCmd.Flags().BoolVar(&resetRegistry, "reset-registry", false, "Clear and rebuild the project registry")
	rootCmd.AddCommand(statusCmd)
}