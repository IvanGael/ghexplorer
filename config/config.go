package config

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	GithubAPIBaseURL = "https://api.github.com"
	HeaderHeight     = 3
	FooterHeight     = 2
)

const (
	ItemsPerPage    = 10
	PaginationStyle = "• %d/%d •"
)

var (
	RepositoryStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ffff"))
	FolderStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("33"))
	FileStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF80"))
	SelectedStyle   = lipgloss.NewStyle().Background(lipgloss.Color("205")).Foreground(lipgloss.Color("#00000"))
	ErrorStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF3333"))
	TabStyle        = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), true, true, false, true).Padding(0, 1)
	ActiveTabStyle  = TabStyle.Border(lipgloss.DoubleBorder(), true, true, false, true)
	SpinnerStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
)

var UseHighPerformanceRenderer = false

var (
	// DocStyle Layout styles
	DocStyle = lipgloss.NewStyle().Padding(1, 2)

	// CardStyle Card styles
	CardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1)

	// ProfileCardStyle Profile card styles
	ProfileCardStyle = CardStyle.BorderForeground(lipgloss.Color("87"))

	// HeaderStyle Header styles
	HeaderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("87")).
			Bold(true).
			Padding(0, 1)

	// FooterStyle Footer styles
	FooterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	// PaginationInfoStyle Pagination styles
	PaginationInfoStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("241")).
				Align(lipgloss.Center)

	// LabelStyle Label styles
	LabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Width(10)

	// ValueStyle Valus styles
	ValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255"))

	// README styles
	// readmeStyle = lipgloss.NewStyle().
	// 		Border(lipgloss.RoundedBorder()).
	// 		BorderForeground(lipgloss.Color("63")).
	// 		Padding(1).
	// 		Width(60)
)
