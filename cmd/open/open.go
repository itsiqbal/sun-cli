/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package open

import (
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/spf13/cobra"
)

var (
	dir string
)

func calculateBasePath(dir string) string{
	result := "~/Desktop/projects/"
	if( dir == "work"){
		result = "~/Desktop/projects/awesomeAirasia/"
	}

	if( dir == "iqbal"){
		result = "~/Desktop/projects/awesomeIqbal/"
	}

	return result;
}

// ListDirectories lists directories in the given path
func ListDirectories(path string) ([]string, error) {
    var directories []string
    files, err := ioutil.ReadDir(path)
    if err != nil {
        return nil, err
    }
    for _, file := range files {
        if file.IsDir() {
            directories = append(directories, file.Name())
        }
    }
    return directories, nil
}

// OpenInVSCode opens the specified directory in VSCode
func OpenInVSCode(dir string) error {
    cmd := exec.Command("code", dir)
    return cmd.Start()
}



// infoCmd represents the info command
var OpenCmd = &cobra.Command{
	Use:   "open",
	Short: "A command pallet related to open directories in vscode",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		basePath := calculateBasePath(dir)

		fmt.Println(basePath)

		// cmd.Help()
	},
}

func init() {

	OpenCmd.Flags().StringVarP(&dir, "directory", "d", "", "set directory to open")

	// OpenCmd.AddCommand(weatherCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
