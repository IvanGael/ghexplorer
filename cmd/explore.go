package cmd

import (
	"fmt"
	"os"

	"ghexplorer/model"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var (
	userFlag   string
	outputFlag string
	formatFlag string
)

func init() {
	exploreCmd := &cobra.Command{
		Use:   "explore [username]",
		Short: "Explore a GitHub profile",
		Long: `Start the TUI application to explore a GitHub profile.
If a username is provided, it will directly load that profile.

Example:
  ghexplorer explore octocat`,
		Run: runExplore,
	}

	// Add flags
	exploreCmd.Flags().StringVarP(&outputFlag, "output", "o", "", "Output file for saving data")
	exploreCmd.Flags().StringVarP(&formatFlag, "format", "f", "text", "Output format (text/json)")

	rootCmd.AddCommand(exploreCmd)
}

func runExplore(cmd *cobra.Command, args []string) {
	var initialGithubID string
	if len(args) > 0 {
		initialGithubID = args[0]
	}

	p := tea.NewProgram(model.InitialModel(initialGithubID), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
