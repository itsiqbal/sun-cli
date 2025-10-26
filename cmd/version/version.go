package version

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	// Version information - injected at build time via ldflags
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// SetVersionInfo sets the version information from main package
// This should be called from main.go with the ldflags values
func SetVersionInfo(v, c, d string) {
	if v != "" {
		version = v
	}
	if c != "" {
		commit = c
	}
	if d != "" {
		date = d
	}
}

// VersionCmd represents the version command
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Long: `Display detailed version information including:
  - Version number
  - Git commit hash
  - Build date and time
  - Go version used for compilation
  - Operating system and architecture`,
	Example: `  # Show version information
  sun version

  # Or use the shorthand flag
  sun --version
  sun -v`,
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}

// GetVersion returns the version string
func GetVersion() string {
	return version
}

// GetCommit returns the git commit hash
func GetCommit() string {
	return commit
}

// GetBuildDate returns the build date
func GetBuildDate() string {
	return date
}

// printVersion displays formatted version information
func printVersion() {
	fmt.Printf("\n")
	fmt.Printf("╭─────────────────────────────────────────╮\n")
	fmt.Printf("│           Sun CLI                       │\n")
	fmt.Printf("╰─────────────────────────────────────────╯\n")
	fmt.Printf("\n")
	fmt.Printf("  Version:      %s\n", version)
	fmt.Printf("  Commit:       %s\n", commit)
	fmt.Printf("  Build Date:   %s\n", date)
	fmt.Printf("  Go Version:   %s\n", runtime.Version())
	fmt.Printf("  OS/Arch:      %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("  Compiler:     %s\n", runtime.Compiler)
	fmt.Printf("\n")
	fmt.Printf("  Repository:   https://github.com/itsiqbal/sun-cli\n")
	fmt.Printf("  Report bugs:  https://github.com/itsiqbal/sun-cli/issues\n")
	fmt.Printf("\n")
}
