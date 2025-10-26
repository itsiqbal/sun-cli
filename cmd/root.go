/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/itsiqbal/sun-cli/cmd/ai"
	"github.com/itsiqbal/sun-cli/cmd/encrypt"
	"github.com/itsiqbal/sun-cli/cmd/gcp"
	"github.com/itsiqbal/sun-cli/cmd/info"
	"github.com/itsiqbal/sun-cli/cmd/version"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sun",
	Short: "Sun CLI - A powerful command-line tool for cloud service management",
	Long: `Sun CLI is a comprehensive command-line interface designed to streamline 
your cloud infrastructure and service management workflow.

With Sun CLI, you can:
  • Quickly open and manage GCP (Google Cloud Platform) services
  • Navigate through cloud resources with ease
  • List and explore directory structures efficiently
  • Automate repetitive cloud operations
  • Integrate with your existing DevOps workflows

Sun CLI is built with performance and developer experience in mind, 
providing fast execution, intuitive commands, and seamless integration 
with modern cloud platforms.

Examples:
  # Open a GCP service in your browser
  sun open compute

  # List directory contents with details
  sun list /path/to/directory

  # View help for any command
  sun [command] --help

For more information and documentation, visit:
  https://github.com/itsiqbal/sun-cli`,
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommand is provided, show help
		_ = cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func addSubCommandPallets() {
	rootCmd.AddCommand(version.VersionCmd)
	rootCmd.AddCommand(info.InfoCmd)
	rootCmd.AddCommand(ai.AiCmd)
	rootCmd.AddCommand(encrypt.EncryptCmd)
	rootCmd.AddCommand(gcp.GcpCmd)

}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.sun-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	addSubCommandPallets()

	// Handle --version flag before command execution
	rootCmd.PreRun = func(cmd *cobra.Command, args []string) {
		if v, _ := cmd.Flags().GetBool("version"); v {
			fmt.Printf("sun version %s\n", version.GetVersion())
			os.Exit(0)
		}
	}
}
