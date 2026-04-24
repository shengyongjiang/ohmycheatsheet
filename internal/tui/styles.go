package tui

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("15")).
			Padding(0, 1)

	headerStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1)

	descStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15"))

	cmdStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("14")).
			Padding(0, 2)

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	reviewTagStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("11"))

	rememberedTagStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	statusBarStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).
			BorderForeground(lipgloss.Color("240")).
			Foreground(lipgloss.Color("240")).
			Padding(0, 1)

	counterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("14"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	selectedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("236")).
			Foreground(lipgloss.Color("15")).
			Bold(true).
			Padding(0, 1)

	selectedCmdStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("236")).
				Foreground(lipgloss.Color("14")).
				Padding(0, 2)

	backfillPageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("13")).
				Italic(true)

	separatorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))
)
