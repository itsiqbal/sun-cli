/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package open

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	dir string
)

func calculateBasePath(dir string) string{
	result := "~/Desktop/projects/"
	if( dir == "work"){
		result = "/Users/aauser/Desktop/projects/awesomeAirasia"
	}

	if( dir == "iqbal"){
		result = "/Users/aauser/Desktop/projects/awesomeIqbal/"
	}

	return result;
}

// ListDirectories lists directories in the given path
func listDirectories(path string) ([]string, error) {
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
func openInVSCode(dir string) error {
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

		directories, err := listDirectories(basePath)

		if err != nil {
			log.Fatal(err)
		}

		if len(directories) == 0 {
			fmt.Println("No directories found.")
			return
		}

		for {
			fmt.Println("Select a directory to open in VSCode:")
			for i, dir := range directories {
				fmt.Printf("%d => %s\n", i+1, dir)
			}
			
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter the number of the directory: ")
			input, _ := reader.ReadString('\n')
			selection, err := strconv.Atoi(input[:len(input)-1])
			if err != nil || selection < 1 || selection > len(directories) {
				if(selection == 99){
					break
				}else{
					fmt.Println("Invalid selection")
				}
				return
			}
			selectedDir := directories[selection-1]
			fmt.Println(selectedDir)
			// err = openInVSCode(filepath.Join(path, selectedDir))
			// if err != nil {
			// 	log.Fatal(err)
			// }
			// fmt.Printf("Opening %s in VSCode...\n", selectedDir)
			

			// if number == 99 {     // Condition to exit the loop
			// 	break
			// }
		}
		

		cmd.Help()
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
