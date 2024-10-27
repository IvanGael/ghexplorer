package main

import (
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles the CLI view interaction updates
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
				switch m.activeTab {
				case 0:
					m.currentView = "profile"
				case 1:
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
		case "up", "down", "pgup", "pgdown":
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
		width, height := m.calculateViewportDimensions(msg)
		m.viewport = viewport.New(width, height)
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

// fetchProfile handles the profile fetching
func (m model) fetchProfile() tea.Msg {
	profile, err := fetchGitHubProfile(m.githubID)
	if err != nil {
		return err
	}

	return *profile
}

// fetchRepositories handles the profile repositories fetching
func (m model) fetchRepositories() tea.Msg {
	repos, err := fetchRepositories(m.profile.Login)
	if err != nil {
		return err
	}
	return repos
}

// fetchRepositoryContents handles the profile repository contents fetching
func (m model) fetchRepositoryContents() tea.Msg {
	contents, err := fetchRepositoryContents(m.profile.Login, m.selected["repository"], m.selected["path"])
	if err != nil {
		return err
	}
	return contents
}

// fetchFileContent handles the profile repository file content fetching
func (m model) fetchFileContent() tea.Msg {
	content, err := fetchFileContent(m.profile.Login, m.selected["repository"], m.selected["path"]+"/"+m.selected["file"])
	if err != nil {
		return err
	}
	return content
}

// searchRepositories handles the profile repositories search performing
func (m model) searchRepositories() tea.Msg {
	repos, err := searchRepositories(m.profile.Login, m.searchQuery)
	if err != nil {
		return err
	}
	return repos
}
