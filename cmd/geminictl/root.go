package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var testbedDir string

var rootCmd = &cobra.Command{
	Use:   "geminictl",
	Short: "geminictl is a session manager for Gemini CLI",
	Long:  `A CLI utility designed to provide observability and management capabilities for Gemini CLI sessions and projects.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Default action if no command is specified
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&testbedDir, "testbed", "", "Path to a testbed directory (overrides default storage)")
}