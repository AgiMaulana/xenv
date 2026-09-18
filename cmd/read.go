package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Read and redact all environment variables",
	Run: func(cmd *cobra.Command, args []string) {
		envVars, err := parseEnv(inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		redacted := make(map[string]string)
		for k := range envVars {
			redacted[k] = "<redacted>"
		}

		if jsonOutput {
			out, _ := json.Marshal(redacted)
			fmt.Println(string(out))
		} else {
			for k, v := range redacted {
				fmt.Printf("%s=%s\n", k, v)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(readCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// readCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// readCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
