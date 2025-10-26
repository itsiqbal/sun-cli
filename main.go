/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/itsiqbal/sun-cli/cmd"
	"github.com/itsiqbal/sun-cli/cmd/version"
)

var (
	versionInfo = "dev"
	commit      = "none"
	date        = "unknown"
)

func main() {
	// Set version info in the version package
	version.SetVersionInfo(versionInfo, commit, date)

	//execute the root command
	cmd.Execute()
}
