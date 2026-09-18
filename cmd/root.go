package cmd

import (
	"bufio"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	inputFile  string
	jsonOutput bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "xenv",
	Short: "AI optimized secret reader and redactor",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&inputFile, "input", "i", "secret", "Path to the secret file")
	rootCmd.PersistentFlags().BoolVarP(&jsonOutput, "json", "j", false, "Output in AI-friendly JSON format")
}

func parseEnv(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	envVars := make(map[string]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			envVars[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	return envVars, scanner.Err()
}
