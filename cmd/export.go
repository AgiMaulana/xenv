/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	exportKey string
	deriveKey string
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Safely generate an export command for eval",
	Run: func(cmd *cobra.Command, args []string) {
		envVars, err := parseEnv(inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}

		val, exists := envVars[deriveKey]
		if !exists {
			fmt.Fprintf(os.Stderr, "Error: key %s not found in %s", deriveKey, inputFile)
			os.Exit(1)
		}

		fileInfo, _ := os.Stdout.Stat()
		if (fileInfo.Mode() & os.ModeCharDevice) != 0 {
			fmt.Fprintln(os.Stderr, "🔒 Security block: Refusing to print secret to the screen.")
			fmt.Fprintln(os.Stderr, "To set the variable, wrap this command in eval:")
			fmt.Fprintf(os.Stderr, "  eval $(./xenv export -i %s -d %s -k %s)\n", inputFile, deriveKey, exportKey)
			os.Exit(1)
		}
		fmt.Printf("export %s=%s\n", exportKey, shellQuote(val))
	},
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
}

func init() {
	rootCmd.AddCommand(exportCmd)

	exportCmd.Flags().StringVarP(&exportKey, "key", "k", "", "The new environment variable")
	exportCmd.Flags().StringVarP(&deriveKey, "derive", "d", "", "The key from the secret file")

	exportCmd.MarkFlagRequired("key")
	exportCmd.MarkFlagRequired("derive")
}
