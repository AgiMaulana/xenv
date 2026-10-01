package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	stateNotExist  = "not-exist"
	stateAvailable = "available"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check [key]",
	Short: "Check if a specific key exists",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		keyToCheck := args[0]

		manager, err := newSourceManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		_, exists := manager.Lookup(keyToCheck)
		state := stateNotExist
		if exists {
			state = stateAvailable
		}

		if jsonOutput {
			printJSON(map[string]interface{}{
				"key":    keyToCheck,
				"exists": exists,
				"state":  state,
			})
		} else {
			fmt.Printf("%s=<%s>\n", keyToCheck, state)
		}

		if !exists {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
