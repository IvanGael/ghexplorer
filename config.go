package main

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	githubAPIBaseURL = "https://api.github.com"
	headerHeight     = 3
	footerHeight     = 2
)

const (
	itemsPerPage    = 10
	paginationStyle = "• %d/%d •"
)

var (
	repositoryStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ffff"))
	folderStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("33"))
	fileStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF80"))
	selectedStyle   = lipgloss.NewStyle().Background(lipgloss.Color("205")).Foreground(lipgloss.Color("#00000"))
	errorStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF3333"))
	tabStyle        = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), true, true, false, true).Padding(0, 1)
	activeTabStyle  = tabStyle.Border(lipgloss.DoubleBorder(), true, true, false, true)
)

var useHighPerformanceRenderer = false

var (
	// Layout styles
	docStyle = lipgloss.NewStyle().Padding(1, 2)

	// Card styles
	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1)

	// Profile card styles
	profileCardStyle = cardStyle.BorderForeground(lipgloss.Color("87"))

	// Header styles
	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("87")).
			Bold(true).
			Padding(0, 1)

	// Footer styles
	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	// Pagination styles
	paginationInfoStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("241")).
				Align(lipgloss.Center)

	// Label styles
	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Width(10)

	// Value styles
	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255"))

	// README styles
	// readmeStyle = lipgloss.NewStyle().
	// 		Border(lipgloss.RoundedBorder()).
	// 		BorderForeground(lipgloss.Color("63")).
	// 		Padding(1).
	// 		Width(60)
)
