/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package ai

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

var (
	prompt string
)

// infoCmd represents the info command
var AiCmd = &cobra.Command{
	Use:   "ai",
	Short: "A command pallet related to ai searches",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(prompt)
		url := "http://localhost:11434/api/chat"
		jsonData := []byte(fmt.Sprintf(`{
		"model": "llama3",
		"messages": [
			{
				"role": "user",
				"content": "%s & Reply should be in 25 words only"
			}
		],
		"stream": false
	}`, prompt))

		// Create a new HTTP request
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		// Set the Content-Type header to application/json
		req.Header.Set("Content-Type", "application/json")

		// Create an HTTP client and set a timeout
		client := &http.Client{Timeout: time.Second * 10}

		// Send the HTTP request
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			return
		}
		// defer resp.Body.Close()

		// Read and print the response
		fmt.Printf("Response status: %s\n", resp.Status)
		buf := new(bytes.Buffer)
		if _, err := buf.ReadFrom(resp.Body); err != nil {
			fmt.Printf("⚠️ Failed to read response body: %v\n", err)
		}
		fmt.Printf("Response body: %s\n", buf.String())
	},
}

func init() {

	AiCmd.Flags().StringVarP(&prompt, "prompt", "p", "", "write prompt to search")
	fmt.Print(prompt)

	// AiCmd.AddCommand(weatherCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
