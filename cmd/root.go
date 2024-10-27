package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ghexplorer",
	Short: "GitHub Explorer - A TUI tool to explore GitHub profiles",
	Long: `GitHub Explorer is a Terminal User Interface (TUI) application that allows you to
explore GitHub profiles, repositories, and files in an interactive way.
Complete documentation is available at https://github.com/IvanGael/ghexplorer`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
