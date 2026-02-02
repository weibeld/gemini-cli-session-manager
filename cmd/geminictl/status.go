package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the status of Gemini CLI projects and sessions",
	Run: func(cmd *cobra.Command, args []string) {
		// This will be implemented in the next tasks to launch the Bubbletea TUI
		fmt.Println("Launching TUI status view...")
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
