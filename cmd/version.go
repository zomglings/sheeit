package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zomglings/sheeit/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		if format == "json" {
			output := map[string]string{"version": version.Version}
			enc := json.NewEncoder(os.Stdout)
			if err := enc.Encode(output); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		} else {
			fmt.Printf("sheeit version %s\n", version.Version)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
