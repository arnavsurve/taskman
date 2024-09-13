package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "taskman",
	Short: "Task management at terminal velocity.",
	Long: `
Task management at terminal velocity.

    Made with <3 by Arnav Surve.
    . arnav@surve.dev
    .. surve.dev
    ... github.com/arnavsurve
    .... linkedin.com/in/arnavsurve`,
}

func init() {
	RootCmd.AddCommand(signupCmd)
	// Add other commands here as you create them

}
