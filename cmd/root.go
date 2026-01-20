package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var format string

var rootCmd = &cobra.Command{
	Use:   "sheeit",
	Short: "A CLI for driving spreadsheets",
	Long:  "sheeit is a command-line tool for driving spreadsheets from the terminal.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&format, "format", "plain", "Output format (plain|json)")
}
