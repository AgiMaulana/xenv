package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Read and redact all environment variables",
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := newSourceManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		redacted := make(map[string]string)
		for key := range manager.All() {
			redacted[key] = "<redacted>"
		}

		if jsonOutput {
			printJSON(redacted)
		} else {
			for key, value := range redacted {
				fmt.Printf("%s=%s\n", key, value)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(readCmd)
}
