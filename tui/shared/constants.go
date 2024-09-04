package shared

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

var (
	FocusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	BlurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	CursorStyle         = FocusedStyle
	NoStyle             = lipgloss.NewStyle()
	HelpStyle           = BlurredStyle
	CursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	FocusedButton = FocusedStyle.Render("[ Submit ]")
	BlurredButton = fmt.Sprintf("[ %s ]", BlurredStyle.Render("Submit"))
)

const URL = "http://localhost:8080"
