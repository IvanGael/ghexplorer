# GitHub Profile Explorer CLI

GitHub Profile Explorer is a terminal-based application written in Go that allows users to interactively explore GitHub profiles, repositories, and file contents. This tool provides a user-friendly interface to navigate through GitHub Profile without leaving your terminal.

https://github.com/user-attachments/assets/6aa06a4f-628c-4f5c-8491-8c3fd35ba060

## Features

- **Profile Viewing**: Enter a GitHub username to view basic profile information.
- **Repository Listing**: Browse through a user's repositories with descriptions.
- **File Navigation**: Explore repository contents, including folders and files.
- **File Content Display**: View the contents of files directly in the terminal.
- **Repository Search**: Search for specific repositories within a user's profile.
- **Interactive Navigation**: Use keyboard shortcuts to navigate through different views.
- **Color-Coded Display**: Repositories, folders, and files are color-coded for easy identification.
- **Scrollable File Content**: Navigate through long file contents using scroll functionality.
- **Text Selection and Copying**: Select and copy file contents to your clipboard.

## Prerequisites

Before you begin, ensure you have the following installed:
- Go (version 1.16 or later)
- Git

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/IvanGael/Git_CLI_Explorer.git
   cd Git_CLI_Explorer
   ```

2. Install the required dependencies:
   ```
   go get github.com/charmbracelet/bubbletea
   go get github.com/charmbracelet/lipgloss
   go get github.com/atotto/clipboard
   ```

3. Build the application:
   ```
   go build
   ```

## Usage

1. Run the application:
   ```
   ./github-profile-explorer
   ```

2. Enter a GitHub username when prompted.

3. Use the following keyboard shortcuts to navigate:
   - Arrow keys: Move cursor / Scroll file contents
   - Enter: Select / Open
   - Esc: Go back / Exit selection mode
   - '/': Enter search mode (when viewing repositories)
   - Ctrl+A: Select all (in file view)
   - Ctrl+C: Copy selected text (in file view)
   - Ctrl+D: Deselect all (in file view)
   - PgUp/PgDown: Scroll file contents quickly
   - 'q': Quit the application

## Customization

You can make any customization regarding styles or colors used in the application by modifying the `config.go` file:

```go
var (
	repositoryStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	folderStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("33"))
	fileStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
	selectedStyle   = lipgloss.NewStyle().Background(lipgloss.Color("25"))
)
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) for the TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) for terminal styling
- [Clipboard](https://github.com/atotto/clipboard) for clipboard functionality

## Disclaimer

This application uses the GitHub API without authentication, which has rate limits. For a production application, consider implementing proper authentication using GitHub tokens to increase the rate limits and access private repositories if needed.