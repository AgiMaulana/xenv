package cmd

import (
	"encoding/json"
	"errors"
	"os"

	"xenv/source"

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
	rootCmd.PersistentFlags().StringVarP(&inputFile, "input", "i", "", "Path to the secret file")
	rootCmd.PersistentFlags().BoolVarP(&jsonOutput, "json", "j", false, "Output in AI-friendly JSON format")
}

// newSourceManager builds the source chain shared by every command. An input
// file, when provided and readable, is used as the source; otherwise the
// process environment is used. A missing file also falls back to the
// environment instead of failing.
func newSourceManager() (*source.Manager, error) {
	var sources []source.Source

	if inputFile != "" {
		fileSource, err := source.NewFileSource(inputFile)
		if err == nil {
			sources = append(sources, fileSource)
		} else if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
	}

	if len(sources) == 0 {
		sources = append(sources, &source.SystemSource{})
	}

	return source.NewManager(sources...), nil
}

// printJSON writes v as JSON without HTML-escaping, so characters like < and >
// used by the state labels stay readable.
func printJSON(v interface{}) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetEscapeHTML(false)
	_ = encoder.Encode(v)
}
