package cmd

import (
	"fmt"

	"github.com/arnavsurve/taskman/tui/views"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var signupCmd = &cobra.Command{
	Use:   "signup",
	Short: "Start the signup process",
	Run: func(cmd *cobra.Command, args []string) {
		if err := startSignupTUI(); err != nil {
			fmt.Printf("Error running signup: %v\n", err)
		}
	},
}

func startSignupTUI() error {
	p := tea.NewProgram(views.NewSignUpModel())
	_, err := p.Run()
	return err
}
