package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

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

By default it downloads the latest GitHub release binary for your
platform and installs it to the same location as the current binary.

Use --version to pin a specific release, e.g.:
  xenv update --version v0.0.1
`,
	Run: func(cmd *cobra.Command, args []string) {
		pinVersion, _ := cmd.Flags().GetString("version")
		update(pinVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().String("version", "", "Pin a specific version (e.g. v0.0.1)")
}

func update(pinVersion string) {
	// Step 1: where does the current binary live?
	currentBinary, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: cannot determine current executable path: %v\n", err)
		os.Exit(1)
	}

	// Resolve symlinks (homebrew, asdf, etc.)
	realBinary, err := filepath.EvalSymlinks(currentBinary)
	if err == nil {
		currentBinary = realBinary
	}

	// Step 2: resolve version
	releaseVersion := pinVersion
	if releaseVersion == "" {
		fmt.Println("Looking up the latest release...")
		out, err := exec.Command("curl", "-fsSL", "-o",
			"/dev/null", "-w", "%{url_effective}",
			"https://github.com/AgiMaulana/xenv/releases/latest",
		).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: could not determine latest version: %v\n", err)
			os.Exit(1)
		}
		releaseVersion = extractVersionFromURL(string(out))
		if releaseVersion == "" {
			fmt.Fprintf(os.Stderr, "Error: could not parse latest release version\n")
			os.Exit(1)
		}
	}

	if !strings.HasPrefix(releaseVersion, "v") {
		releaseVersion = "v" + releaseVersion
	}

	// Already up to date?
	if releaseVersion == version && pinVersion == "" {
		fmt.Printf("Already at the latest version (%s).\n", version)
		return
	}

	// Step 3: build download URL
	asset := fmt.Sprintf("xenv-%s-%s", runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		asset += ".exe"
	}

	downloadURL := fmt.Sprintf(
		"https://github.com/AgiMaulana/xenv/releases/download/%s/%s",
		releaseVersion, asset,
	)

	fmt.Printf("Downloading %s...\n", asset)

	// Step 4: download to temp file
	tmpFile, err := os.CreateTemp("", "xenv-*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not create temp file: %v\n", err)
		os.Exit(1)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)

	download := exec.Command("curl", "-fsSL", "--retry", "3", "-o", tmpPath, downloadURL)
	download.Stderr = os.Stderr
	if err := download.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "\nDownload failed. Does release %s contain %s?\n", releaseVersion, asset)
		os.Exit(1)
	}

	// Step 5: make executable
	if err := os.Chmod(tmpPath, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not make binary executable: %v\n", err)
		os.Exit(1)
	}

	// Step 6: atomically replace current binary
	// Rename current -> backup, then move new -> current
	backupPath := currentBinary + ".bak"
	os.Remove(backupPath) // clean up any stale backup

	if err := os.Rename(currentBinary, backupPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not backup current binary: %v\n", err)
		os.Exit(1)
	}

	if err := os.Rename(tmpPath, currentBinary); err != nil {
		// Restore backup
		os.Rename(backupPath, currentBinary)
		fmt.Fprintf(os.Stderr, "Error: could not install update: %v\n", err)
		os.Exit(1)
	}

	os.Remove(backupPath) // success, clean up

	fmt.Printf("\n✓ xenv updated %s → %s\n", version, releaseVersion)
	fmt.Printf("  Installed at: %s\n", currentBinary)
}

func extractVersionFromURL(url string) string {
	// URL: https://github.com/AgiMaulana/xenv/releases/tag/v0.0.1
	const prefix = "/tag/"
	idx := strings.LastIndex(url, prefix)
	if idx < 0 {
		return ""
	}
	return url[idx+len(prefix):]
}
