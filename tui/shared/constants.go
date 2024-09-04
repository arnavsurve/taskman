package shared

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

var (
	// FocusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	// FocusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#B394CB"))
	FocusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575"))

	BlurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	CursorStyle         = FocusedStyle
	HelpStyle           = BlurredStyle
	NoStyle             = lipgloss.NewStyle()
	CursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	HeaderStyle         = lipgloss.NewStyle().Bold(true)

	FocusedButton = FocusedStyle.Render("[ Submit ]")
	BlurredButton = fmt.Sprintf("[ %s ]", BlurredStyle.Render("Submit"))
)

const URL = "http://localhost:8080"
