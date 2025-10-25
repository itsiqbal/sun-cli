package version

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	// Version information - injected at build time via ldflags
	// Build with: go build -ldflags="-X cmd.version=1.0.0"
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// VersionCmd represents the info command
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

func init() {
	VersionCmd.Flags().BoolP("version", "v", false, "Print version information")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
