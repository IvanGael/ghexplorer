package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"ghexplorer/github_api"
	"github.com/spf13/cobra"
)

func init() {
	searchCmd := &cobra.Command{
		Use:   "search [username] [query]",
		Short: "Search repositories for a user",
		Long: `Search through a user's repositories with a query string.
This command provides search results without starting the TUI.

Example:
  ghexplorer search octocat "awesome"`,
		Args: cobra.ExactArgs(2),
		Run:  runSearch,
	}

	// Add flags
	searchCmd.Flags().StringVarP(&outputFlag, "output", "o", "", "Output file for saving data")
	searchCmd.Flags().StringVarP(&formatFlag, "format", "f", "text", "Output format (text/json)")

	rootCmd.AddCommand(searchCmd)
}

func runSearch(cmd *cobra.Command, args []string) {
	username := args[0]
	query := args[1]

	// Search repositories
	repos, err := github_api.SearchRepositories(username, query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Handle output based on format flag
	switch formatFlag {
	case "json":
		output, err := json.MarshalIndent(repos, "", "  ")
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
			for _, repo := range repos {
				fmt.Fprintf(f, "Repository: %s\nDescription: %s\n\n", repo.Name, repo.Description)
			}
		} else {
			for _, repo := range repos {
				fmt.Printf("Repository: %s\nDescription: %s\n\n", repo.Name, repo.Description)
			}
		}
	}
}
