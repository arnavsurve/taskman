package shared

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

var (
	FrenchGray = "#BBBCC5" // kinda light purple
	White      = "#FCFCFC"
	Verdigris  = "#5B9A96" // teal
	DimGray    = "#677176"
	Wenge      = "#5A5252" // brown

	// FocusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	// FocusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#B394CB"))
	FocusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(Verdigris))

	BlurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	CursorStyle         = FocusedStyle
	HelpStyle           = BlurredStyle
	NoStyle             = lipgloss.NewStyle()
	CursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	HeaderStyle         = lipgloss.NewStyle().Bold(true)

	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(White)).
			Background(lipgloss.Color(Verdigris)).
			Padding(0, 1).
			Margin(1)

	FocusedButton = FocusedStyle.Render("[ Submit ]")
	BlurredButton = fmt.Sprintf("[ %s ]", BlurredStyle.Render("Submit"))
)

const URL = "http://localhost:8080"
