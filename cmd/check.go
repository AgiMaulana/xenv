package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check [key]",
	Short: "Check if a specific key exists",
	Run: func(cmd *cobra.Command, args []string) {
		keyToCheck := args[0]
		envVars, err := parseEnv(inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		_, exists := envVars[keyToCheck]

		if jsonOutput {
			out, _ := json.Marshal(map[string]interface{}{"key": keyToCheck, "exists": exists})
			fmt.Println(string(out))
		} else {
			status := "<missing>"
			if exists {
				status = "<available>"
			}
			fmt.Printf("%s=%s\n", keyToCheck, status)
		}

		if !exists {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
