package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"geminictl/internal/cache"
	"geminictl/internal/scanner"
	"geminictl/internal/tui"
)

var resetRegistry bool

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the status of Gemini CLI projects and sessions",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := cache.NewCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing cache: %v\n", err)
			os.Exit(1)
		}

		if resetRegistry {
			c.Clear()
			_ = c.Save()
		} else {
			if err := c.Load(); err != nil {
				fmt.Fprintf(os.Stderr, "Error loading cache: %v\n", err)
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

		// --- Integrity Check: Garbage Collection ---
		// Remove cache entries that are no longer present in ~/.gemini/tmp
		activeIDs := make(map[string]bool)
		for _, p := range projects {
			activeIDs[p.ID] = true
		}

		changed := false
		for hash := range c.Data {
			if !activeIDs[hash] {
				c.Delete(hash)
				changed = true
			}
		}
		if changed {
			_ = c.Save()
		}

		m := tui.NewModel(projects, c, scan)
		p := tea.NewProgram(m, tea.WithAltScreen())

		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	statusCmd.Flags().BoolVar(&resetRegistry, "reset-registry", false, "Clear and rebuild the project cache")
	rootCmd.AddCommand(statusCmd)
}
