package main

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// GitHubProfile is github profile struct
type GitHubProfile struct {
	Name        string `json:"name"`
	Login       string `json:"login"`
	Description string `json:"bio"`
	Followers   int    `json:"followers"`
	Following   int    `json:"following"`
}

// Repository is github profile respository struct
type Repository struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// FileInfo is github profile repository file info struct
type FileInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// model is the base model
type model struct {
	githubID     string
	inputting    bool
	profile      *GitHubProfile
	repositories []*Repository
	readme       string
	currentView  string
	cursor       int
	selected     map[string]string
	fileContents []*FileInfo
	fileContent  string
	searchQuery  string
	errorMessage string
	selectMode   bool
	selectStart  int
	selectEnd    int
	textInput    textinput.Model
	viewport     viewport.Model
	spinner      spinner.Model
	tabs         []string
	activeTab    int
}

// initialModel initialize the model
func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter GitHub profile ID"
	ti.Focus()

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	vp := viewport.New(80, 20) // Default size, will be adjusted later
	vp.YPosition = headerHeight

	return model{
		githubID:    "",
		inputting:   true,
		currentView: "input",
		selected:    make(map[string]string),
		textInput:   ti,
		spinner:     s,
		tabs:        []string{"Overview", "Repositories"},
		activeTab:   0,
		viewport:    vp,
	}
}

// Init init the model
func (m model) Init() tea.Cmd {
	return textinput.Blink
}
