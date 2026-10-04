package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/spf13/cobra"
)

// injectCmd represents the inject command
var injectCmd = &cobra.Command{
	Use:   "inject KEY... [--] COMMAND [ARG...]",
	Short: "Inject secrets into a command's environment and execute it",
	Long: `Resolve one or more secrets and pass them as environment variables
to the given command, then execute it.

The secret values are never printed to the terminal. They exist only
in the child process's environment for the duration of the command.

Examples:
  xenv inject API_KEY -- sh -c 'curl -H "Authorization: Bearer $API_KEY" https://api.example.com'
  xenv inject GITHUB_TOKEN AWS_ACCESS_KEY_ID -- ./deploy.sh

Wrap commands that reference the injected variable in sh -c '...' (or a
script): your shell expands $VAR before xenv runs, so a bare command line
expands to the parent's value, not the injected one.

If a requested key is not found in any source, inject fails with an
error before starting the command.
`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := newSourceManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Keys come before the "--" separator; the command follows it.
		// pflag consumes the literal "--" before Run sees args, so we recover
		// its position with ArgsLenAtDash. Without a separator, only the first
		// arg is a key and the rest is the command.
		var keys []string
		var commandArgs []string

		if dash := cmd.Flags().ArgsLenAtDash(); dash >= 0 {
			keys = args[:dash]
			commandArgs = args[dash:]
		} else {
			keys = args[:1]
			commandArgs = args[1:]
		}

		if len(keys) == 0 {
			fmt.Fprintf(os.Stderr, "Error: no keys specified\n")
			os.Exit(1)
		}

		if len(commandArgs) == 0 {
			fmt.Fprintf(os.Stderr, "Error: no command specified\n")
			os.Exit(1)
		}

		// Resolve all secrets before starting the child process.
		// Fail fast if any key is missing.
		injectEnv := make(map[string]string)
		for _, key := range keys {
			val, found := manager.Lookup(key)
			if !found {
				fmt.Fprintf(os.Stderr, "Error: key %q not found in any source\n", key)
				os.Exit(1)
			}
			injectEnv[key] = val
		}

		// Build command with the injected environment.
		binary, err := exec.LookPath(commandArgs[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Create the child's environment: inherit current env, then overlay
		// injected variables.
		env := os.Environ()
		for k, v := range injectEnv {
			env = append(env, k+"="+v)
		}

		// Use syscall.Exec on unix to replace the xenv process with the child,
		// so the child inherits xenv's pid and signal handling. This means
		// the injected env vars are set before the child starts, and when
		// the child exits, xenv is already gone — no extra wait step.
		err = syscall.Exec(binary, commandArgs, env)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to exec %s: %v\n", commandArgs[0], err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(injectCmd)
}
