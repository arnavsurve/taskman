package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "taskman",
	Short: "A task management CLI application",
	Long:  `A task management CLI application built with Go, Cobra, and Bubble Tea.`,
}

func init() {
	RootCmd.AddCommand(signupCmd)
	// Add other commands here as you create them
}
