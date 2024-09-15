package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

type model struct {
	githubID     string
	inputting    bool
	profile      *GitHubProfile
	repositories []*Repository
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

var (
	repositoryStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ffff"))
	folderStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("33"))
	fileStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF80"))
	selectedStyle   = lipgloss.NewStyle().Background(lipgloss.Color("205")).Foreground(lipgloss.Color("#ffffff"))
	errorStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF3333"))
	tabStyle        = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), true, true, false, true).Padding(0, 1)
	activeTabStyle  = tabStyle.Border(lipgloss.DoubleBorder(), true, true, false, true)
	headerHeight    = 3
	footerHeight    = 2
)

var useHighPerformanceRenderer = false

type GitHubProfile struct {
	Name        string `json:"name"`
	Login       string `json:"login"`
	Description string `json:"bio"`
	Followers   int    `json:"followers"`
	Following   int    `json:"following"`
}

type Repository struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type FileInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter GitHub profile ID"
	ti.Focus()

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		githubID:    "",
		inputting:   true,
		currentView: "input",
		selected:    make(map[string]string),
		textInput:   ti,
		spinner:     s,
		tabs:        []string{"Overview", "Repositories"},
		activeTab:   0,
		viewport:    viewport.New(80, 20),
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "tab":
			if m.currentView == "profile" || m.currentView == "repositories" {
				m.activeTab = (m.activeTab + 1) % len(m.tabs)
				if m.activeTab == 0 {
					m.currentView = "profile"
				} else {
					m.currentView = "repositories"
				}
			}
		case "enter":
			switch m.currentView {
			case "input":
				m.inputting = false
				m.currentView = "profile"
				return m, m.fetchProfile
			case "repositories":
				if m.cursor < len(m.repositories) {
					m.selected["repository"] = m.repositories[m.cursor].Name
					m.currentView = "files"
					m.cursor = 0
					m.selected["path"] = ""
					return m, m.fetchRepositoryContents
				}
			case "files":
				if m.cursor < len(m.fileContents) {
					if m.fileContents[m.cursor].Type == "file" {
						m.selected["file"] = m.fileContents[m.cursor].Name
						m.currentView = "fileContent"
						return m, m.fetchFileContent
					} else {
						m.selected["path"] += "/" + m.fileContents[m.cursor].Name
						return m, m.fetchRepositoryContents
					}
				}
			case "search":
				m.currentView = "repositories"
				return m, m.searchRepositories
			}
		case "backspace":
			if m.inputting && len(m.githubID) > 0 {
				m.githubID = m.githubID[:len(m.githubID)-1]
			} else if m.currentView == "search" && len(m.searchQuery) > 0 {
				m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
			}
		case "esc":
			if m.selectMode {
				m.selectMode = false
				m.selectStart = 0
				m.selectEnd = 0
			} else {
				switch m.currentView {
				case "repositories":
					m.currentView = "profile"
				case "files":
					if m.selected["path"] == "" {
						m.currentView = "repositories"
						m.cursor = 0
					} else {
						paths := strings.Split(m.selected["path"], "/")
						m.selected["path"] = strings.Join(paths[:len(paths)-1], "/")
						return m, m.fetchRepositoryContents
					}
				case "fileContent":
					m.currentView = "files"
					m.selectMode = false
					m.selectStart = 0
					m.selectEnd = 0
				case "search":
					m.currentView = "repositories"
				}
			}
		case "up", "down", "left", "right", "pgup", "pgdown":
			if m.currentView == "fileContent" {
				if m.selectMode {
					switch msg.String() {
					case "up":
						m.selectEnd = max(0, m.selectEnd-m.viewport.Width)
					case "down":
						m.selectEnd = min(len(m.fileContent), m.selectEnd+m.viewport.Width)
					case "left":
						m.selectEnd = max(0, m.selectEnd-1)
					case "right":
						m.selectEnd = min(len(m.fileContent), m.selectEnd+1)
					}
				}
				m.viewport, cmd = m.viewport.Update(msg)
				return m, cmd
			} else {
				switch msg.String() {
				case "up":
					if m.cursor > 0 {
						m.cursor--
					}
				case "down":
					switch m.currentView {
					case "repositories":
						if m.cursor < len(m.repositories)-1 {
							m.cursor++
						}
					case "files":
						if m.cursor < len(m.fileContents)-1 {
							m.cursor++
						}
					}
				}
			}
		case "ctrl+a":
			if m.currentView == "fileContent" {
				m.selectMode = true
				m.selectStart = 0
				m.selectEnd = len(m.fileContent)
			}
		case "ctrl+c":
			if m.selectMode && m.currentView == "fileContent" {
				selectedText := m.fileContent[m.selectStart:m.selectEnd]
				clipboard.WriteAll(selectedText)
				m.selectMode = false
				m.selectStart = 0
				m.selectEnd = 0
				return m, nil
			}
		case "ctrl+d":
			if m.selectMode && m.currentView == "fileContent" {
				m.selectMode = false
				m.selectStart = 0
				m.selectEnd = len(m.fileContent)
				return m, nil
			}
		case "/":
			if m.currentView == "repositories" {
				m.currentView = "search"
				m.searchQuery = ""
			}
		default:
			if m.inputting {
				m.githubID += msg.String()
			} else if m.currentView == "search" {
				m.searchQuery += msg.String()
			}
		}
	case tea.WindowSizeMsg:
		m.viewport = viewport.New(msg.Width, msg.Height-headerHeight-footerHeight)
		m.viewport.YPosition = headerHeight
		m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
		if m.currentView == "fileContent" {
			m.viewport.SetContent(m.fileContent)
		}
	case GitHubProfile:
		m.profile = &msg
		return m, m.fetchRepositories
	case []*Repository:
		m.repositories = msg
		m.currentView = "repositories"
		m.cursor = 0
	case []*FileInfo:
		m.fileContents = msg
		m.cursor = 0
	case string:
		m.fileContent = msg
		if m.currentView == "fileContent" {
			m.viewport.SetContent(m.fileContent)
			m.viewport.GotoTop()
		}
	case error:
		m.currentView = "error"
		m.errorMessage = msg.Error()
		return m, nil
	}
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	switch m.currentView {
	case "input":
		return fmt.Sprintf("Git CLI Explorer\n\n%s", m.textInput.View())
	case "profile", "repositories":
		return m.tabView()
	case "files":
		return m.filesView()
	case "fileContent":
		return m.fileContentView()
	case "search":
		return m.searchView()
	case "error":
		return m.errorView()
	default:
		return fmt.Sprintf("%s Loading...", m.spinner.View())
	}
}

func (m model) tabView() string {
	doc := strings.Builder{}

	// Render tabs
	renderedTabs := []string{}
	for i, t := range m.tabs {
		if i == m.activeTab {
			renderedTabs = append(renderedTabs, activeTabStyle.Render(t))
		} else {
			renderedTabs = append(renderedTabs, tabStyle.Render(t))
		}
	}
	doc.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...))
	doc.WriteString("\n\n")

	// Render content based on active tab
	switch m.activeTab {
	case 0: // Overview
		doc.WriteString(m.overviewView())
	case 1: // Repositories
		doc.WriteString(m.repositoriesView())
	}

	return doc.String()
}

func (m model) overviewView() string {
	if m.profile == nil {
		return "Loading profile..."
	}

	content := fmt.Sprintf("Name: %s\nUsername: %s\nBio: %s\nFollowers: %d\nFollowing: %d\n\n",
		stringOrNA(m.profile.Name),
		stringOrNA(m.profile.Login),
		stringOrNA(m.profile.Description),
		m.profile.Followers,
		m.profile.Following)

	m.viewport.SetContent(wordwrap.String(content, m.viewport.Width))
	return m.viewport.View()
}

func (m model) repositoriesView() string {
	var content strings.Builder
	for i, repo := range m.repositories {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		repoName := repositoryStyle.Render(stringOrNA(repo.Name))
		content.WriteString(fmt.Sprintf("%s %s: %s\n", cursor, repoName, stringOrNA(repo.Description)))
	}
	content.WriteString("\nPress Enter to view files, '/' to search, Tab to switch tabs")

	m.viewport.SetContent(content.String())
	return m.viewport.View()
}

func (m model) filesView() string {
	var content strings.Builder
	content.WriteString(fmt.Sprintf("Files in %s%s:\n\n", stringOrNA(m.selected["repository"]), stringOrNA(m.selected["path"])))
	for i, file := range m.fileContents {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		var fileNameStyle lipgloss.Style
		if file.Type == "dir" {
			fileNameStyle = folderStyle
		} else {
			fileNameStyle = fileStyle
		}
		fileName := fileNameStyle.Render(stringOrNA(file.Name))
		content.WriteString(fmt.Sprintf("%s %s (%s)\n", cursor, fileName, stringOrNA(file.Type)))
	}
	content.WriteString("\nPress Enter to view file content, Esc to go back")

	m.viewport.SetContent(content.String())
	return m.viewport.View()
}

func (m model) fileContentView() string {
	// Calculate the available height for the viewport
	viewportHeight := m.viewport.Height - headerHeight - footerHeight

	// Create the header
	header := fmt.Sprintf("File: %s\n\n", stringOrNA(m.selected["file"]))

	// Create the footer
	footer := "\nPress Esc to go back, Ctrl+A to select all, Ctrl+C to copy selection, Ctrl+D to deselect all,"

	// Set the content and size of the viewport
	m.viewport.Height = viewportHeight

	// Apply selection styling if in select mode
	var styledContent string
	if m.selectMode {
		before := m.fileContent[:m.selectStart]
		selected := selectedStyle.Render(m.fileContent[m.selectStart:m.selectEnd])
		after := m.fileContent[m.selectEnd:]
		styledContent = before + selected + after
	} else {
		styledContent = m.fileContent
	}

	m.viewport.SetContent(styledContent)

	// Combine the header, viewport content, and footer
	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		m.viewport.View(),
		footer,
	)
}

func (m model) searchView() string {
	return fmt.Sprintf("Find a repository : %s\n\nPress Enter to search, Esc to cancel", m.searchQuery)
}

func (m model) errorView() string {
	errorText := errorStyle.Render("An error occurred: ")
	return fmt.Sprintf("%s %s\n\nPress any key to go back to input", errorText, m.errorMessage)
	// return fmt.Sprintf("An error occurred: %s\n\nPress any key to go back to input", m.errorMessage)
}

func (m model) fetchProfile() tea.Msg {
	// username := strings.TrimPrefix(m.githubID, "https://github.com/")
	profile, err := fetchGitHubProfile(m.githubID)
	if err != nil {
		return err
	}
	return *profile
}

func (m model) fetchRepositories() tea.Msg {
	repos, err := fetchRepositories(m.profile.Login)
	if err != nil {
		return err
	}
	return repos
}

func (m model) fetchRepositoryContents() tea.Msg {
	contents, err := fetchRepositoryContents(m.profile.Login, m.selected["repository"], m.selected["path"])
	if err != nil {
		return err
	}
	return contents
}

func (m model) fetchFileContent() tea.Msg {
	content, err := fetchFileContent(m.profile.Login, m.selected["repository"], m.selected["path"]+"/"+m.selected["file"])
	if err != nil {
		return err
	}
	return content
}

func (m model) searchRepositories() tea.Msg {
	repos, err := searchRepositories(m.profile.Login, m.searchQuery)
	if err != nil {
		return err
	}
	return repos
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
