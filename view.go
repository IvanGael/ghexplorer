package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

// View handles the CLI global view
func (m model) View() string {
	switch m.currentView {
	case "input":
		return docStyle.Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				headerStyle.Render("Git CLI Explorer"),
				"\n",
				cardStyle.Render(m.textInput.View()),
			),
		)
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
		return docStyle.Render(
			cardStyle.Render(
				lipgloss.JoinHorizontal(
					lipgloss.Center,
					m.spinner.View(),
					" Loading...",
				),
			),
		)
	}
}

// calculateViewportDimensions calculates the available viewport dimensions
func (m model) calculateViewportDimensions(msg tea.WindowSizeMsg) (width, height int) {
	// Subtract padding and borders from total width
	width = msg.Width - 4 // 2 for left padding + 2 for right padding

	// Subtract header, footer, and padding from total height
	height = msg.Height - headerHeight - footerHeight - 2 // 2 for top/bottom padding

	return width, height
}

// getPaginationInfo returns pagination details for the current view
func (m model) getPaginationInfo() (currentPage, totalPages, startIdx, endIdx int) {
	var totalItems int

	switch m.currentView {
	case "repositories":
		totalItems = len(m.repositories)
	case "files":
		totalItems = len(m.fileContents)
	default:
		return 1, 1, 0, 0
	}

	currentPage = (m.cursor / itemsPerPage) + 1
	totalPages = (totalItems + itemsPerPage - 1) / itemsPerPage

	startIdx = (currentPage - 1) * itemsPerPage
	endIdx = min(startIdx+itemsPerPage, totalItems)

	return currentPage, totalPages, startIdx, endIdx
}

// renderPagination renders the pagination information
func renderPagination(current, total int) string {
	if total <= 1 {
		return ""
	}
	return paginationInfoStyle.Render(fmt.Sprintf(paginationStyle, current, total))
}

// tabView handles the CLI tab view
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

	return docStyle.Render(doc.String())
}

// overviewView handles the CLI overview view
func (m model) overviewView() string {
	if m.profile == nil {
		return cardStyle.Render("Loading profile...")
	}

	// Calculate widths for split view
	totalWidth := m.viewport.Width
	profileWidth := totalWidth * 2 / 5           // 40% of width
	readmeWidth := totalWidth - profileWidth - 4 // Remaining width minus padding

	profileCardStyle = profileCardStyle.Width(profileWidth)
	readmeStyle = readmeStyle.Width(readmeWidth)

	// Profile information section
	profileInfo := lipgloss.JoinVertical(
		lipgloss.Left,
		headerStyle.Render("Profile Information"),
		lipgloss.JoinHorizontal(lipgloss.Left, labelStyle.Render("Name:"), valueStyle.Render(stringOrNA(m.profile.Name))),
		lipgloss.JoinHorizontal(lipgloss.Left, labelStyle.Render("Username:"), valueStyle.Render(stringOrNA(m.profile.Login))),
		lipgloss.JoinHorizontal(lipgloss.Left, labelStyle.Render("Bio:"), valueStyle.Render(stringOrNA(m.profile.Description))),
		"",
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			labelStyle.Render("Stats:"),
			valueStyle.Render(
				fmt.Sprintf("👥 %d followers • %d following",
					m.profile.Followers,
					m.profile.Following,
				),
			),
		),
	)

	// README section
	readme := "Loading README..."
	if m.readme != "" {
		readme = m.readme
	}

	readmeSection := lipgloss.JoinVertical(
		lipgloss.Left,
		headerStyle.Render("Profile README"),
		wordwrap.String(readme, readmeStyle.GetWidth()-4), // Account for padding
	)

	// Combine sections horizontally
	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		profileCardStyle.Render(profileInfo),
		readmeStyle.Render(readmeSection),
	)

	return content
}

// repositoriesView handles the CLI repositories view
func (m model) repositoriesView() string {
	var content strings.Builder

	currentPage, totalPages, startIdx, endIdx := m.getPaginationInfo()

	content.WriteString(headerStyle.Render("Repositories"))
	content.WriteString("\n\n")

	// Display only the repositories for the current page
	visibleRepos := m.repositories[startIdx:endIdx]
	for i, repo := range visibleRepos {
		cursor := " "
		if startIdx+i == m.cursor {
			cursor = ">"
		}

		repoCard := lipgloss.JoinVertical(
			lipgloss.Left,
			repositoryStyle.Render(stringOrNA(repo.Name)),
			valueStyle.Render(stringOrNA(repo.Description)),
		)

		if startIdx+i == m.cursor {
			repoCard = selectedStyle.Render(repoCard)
		} else {
			repoCard = cardStyle.Render(repoCard)
		}

		content.WriteString(fmt.Sprintf("%s %s\n", cursor, repoCard))
	}

	// Add pagination info
	content.WriteString("\n")
	content.WriteString(renderPagination(currentPage, totalPages))

	footer := footerStyle.Render("\nPress Enter to view files • '/' to search • Tab to switch tabs • ←/→ to change pages")

	content.WriteString(footer)

	return content.String()
}

// filesView handles the CLI files view
func (m model) filesView() string {
	var content strings.Builder

	currentPage, totalPages, startIdx, endIdx := m.getPaginationInfo()

	header := lipgloss.JoinVertical(
		lipgloss.Left,
		headerStyle.Render(fmt.Sprintf("Repository: %s", stringOrNA(m.selected["repository"]))),
		valueStyle.Render(fmt.Sprintf("Path: %s", stringOrNA(m.selected["path"]))),
	)

	content.WriteString(cardStyle.Render(header))
	content.WriteString("\n\n")

	// Display only the files for the current page
	visibleFiles := m.fileContents[startIdx:endIdx]
	for i, file := range visibleFiles {
		cursor := " "
		if startIdx+i == m.cursor {
			cursor = ">"
		}

		var fileNameStyle lipgloss.Style
		var icon string
		if file.Type == "dir" {
			fileNameStyle = folderStyle
			icon = "📁"
		} else {
			fileNameStyle = fileStyle
			icon = "📄"
		}

		fileCard := lipgloss.JoinHorizontal(
			lipgloss.Left,
			icon,
			" ",
			fileNameStyle.Render(stringOrNA(file.Name)),
		)

		if startIdx+i == m.cursor {
			fileCard = selectedStyle.Render(fileCard)
		} else {
			fileCard = cardStyle.Render(fileCard)
		}

		content.WriteString(fmt.Sprintf("%s %s\n", cursor, fileCard))
	}

	// Add pagination info
	content.WriteString("\n")
	content.WriteString(renderPagination(currentPage, totalPages))

	footer := footerStyle.Render("\nPress Enter to view content • Esc to go back • ←/→ to change pages")

	content.WriteString(footer)

	return content.String()
}

// fileContentView handles the CLI fileContent view
func (m model) fileContentView() string {
	header := cardStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			headerStyle.Render("File Content"),
			valueStyle.Render(stringOrNA(m.selected["file"])),
		),
	)

	footer := footerStyle.Render("\nPress Esc to go back • Ctrl+A to select all • Ctrl+C to copy • Ctrl+D to deselect • ↑/↓ to scroll")

	var styledContent string
	if m.selectMode {
		before := m.fileContent[:m.selectStart]
		selected := selectedStyle.Render(m.fileContent[m.selectStart:m.selectEnd])
		after := m.fileContent[m.selectEnd:]
		styledContent = before + selected + after
		return lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			"\n",
			cardStyle.Render(styledContent),
			footer,
		)
	} else {
		styledContent = m.fileContent
	}

	// Set the viewport content if it hasn't been set
	if m.viewport.Height == 0 {
		m.viewport.Height = m.viewport.Height - headerHeight - footerHeight - 4 // Adjust for header, footer, and padding
		m.viewport.Width = m.viewport.Width - 4                                 // Adjust for left and right padding
		m.viewport.SetContent(styledContent)
	}

	// Render the viewport
	content := m.viewport.View()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"\n",
		cardStyle.Render(content),
		footer,
	)
}

// searchView handles the CLI search view
func (m model) searchView() string {
	return docStyle.Render(
		cardStyle.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				headerStyle.Render("Search Repositories"),
				"",
				valueStyle.Render(m.searchQuery),
				"",
				lipgloss.NewStyle().
					Foreground(lipgloss.Color("241")).
					Render("Press Enter to search • Esc to cancel"),
			),
		),
	)
}

// errorView handles the CLI error view
func (m model) errorView() string {
	return docStyle.Render(
		cardStyle.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				errorStyle.Render("Error Occurred"),
				"",
				valueStyle.Render(m.errorMessage),
				"",
				lipgloss.NewStyle().
					Foreground(lipgloss.Color("241")).
					Render("Press any key to go back"),
			),
		),
	)
}
