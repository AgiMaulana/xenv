package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var version = "dev" // set by ldflags at build time

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the current xenv version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update xenv to the latest version",
	Long: `Download and install the latest xenv release.

Re-runs the official install script, which downloads the latest GitHub
release binary for your platform, verifies its cosign keyless signature
(when cosign is available), and installs it to /usr/local/bin
(or \$XENV_INSTALL_DIR if set).

This is the same as:
  curl -fsSL https://raw.githubusercontent.com/AgiMaulana/xenv/main/install.sh | sh
`,
	Run: func(cmd *cobra.Command, args []string) {
		update()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(updateCmd)
}

func update() {
	fmt.Printf("Current version: %s\n", version)
	fmt.Println("Downloading and running the installer...")

	installScript := "https://raw.githubusercontent.com/AgiMaulana/xenv/main/install.sh"

	curl := exec.Command("curl", "-fsSL", installScript)
	sh := exec.Command("sh")

	// Pipe curl stdout -> sh stdin
	pipe, err := curl.StdoutPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not create pipe: %v\n", err)
		os.Exit(1)
	}
	sh.Stdin = pipe

	// Pass through so sudo prompts work and user sees progress
	sh.Stdout = os.Stdout
	sh.Stderr = os.Stderr
	curl.Stderr = os.Stderr

	if err := sh.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not start sh: %v\n", err)
		os.Exit(1)
	}

	if err := curl.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "\nDownload failed. Check your internet connection.\n")
		os.Exit(1)
	}

	if err := sh.Wait(); err != nil {
		fmt.Fprintf(os.Stderr, "\nInstaller exited with an error.\n")
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("Update complete! Run 'xenv version' to verify.")
}
