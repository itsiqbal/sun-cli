/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/itsiqbal/sun-cli/cmd/ai"
	"github.com/itsiqbal/sun-cli/cmd/encrypt"
	"github.com/itsiqbal/sun-cli/cmd/gcp"
	"github.com/itsiqbal/sun-cli/cmd/info"
	"github.com/itsiqbal/sun-cli/cmd/open"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sun",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

func addSubCommandPallets() {
	rootCmd.AddCommand(info.InfoCmd)
	rootCmd.AddCommand(open.OpenCmd)
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
}
