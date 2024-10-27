package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"ghexplorer/github_api"
	"github.com/spf13/cobra"
)

func init() {
	repoCmd := &cobra.Command{
		Use:   "repo [username] [repository]",
		Short: "Get information about a specific repository",
		Long: `Fetch and display information about a specific GitHub repository.
This command provides detailed information without starting the TUI.

Example:
  ghexplorer repo octocat Hello-World`,
		Args: cobra.ExactArgs(2),
		Run:  runRepo,
	}

	// Add flags
	repoCmd.Flags().StringVarP(&outputFlag, "output", "o", "", "Output file for saving data")
	repoCmd.Flags().StringVarP(&formatFlag, "format", "f", "text", "Output format (text/json)")

	rootCmd.AddCommand(repoCmd)
}

func runRepo(cmd *cobra.Command, args []string) {
	username := args[0]
	repository := args[1]

	// Fetch repository contents
	contents, err := github_api.FetchRepositoryContents(username, repository, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Handle output based on format flag
	switch formatFlag {
	case "json":
		output, err := json.MarshalIndent(contents, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if outputFlag != "" {
			err = os.WriteFile(outputFlag, output, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
				os.Exit(1)
			}
		} else {
			fmt.Println(string(output))
		}
	default:
		if outputFlag != "" {
			f, err := os.Create(outputFlag)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
				os.Exit(1)
			}
			defer f.Close()
			for _, content := range contents {
				fmt.Fprintf(f, "%s (%s)\n", content.Name, content.Type)
			}
		} else {
			for _, content := range contents {
				fmt.Printf("%s (%s)\n", content.Name, content.Type)
			}
		}
	}
}
