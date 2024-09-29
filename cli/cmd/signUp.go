/*
Copyright © 2024 Arnav Surve arnav@surve.dev>
*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var github bool

// signUpCmd represents the signUp command
var signUpCmd = &cobra.Command{
	Use:   "signup",
	Short: "Register with taskman.",
	Long:  `Register with taskman.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("signup called")
		if github {
			fmt.Println("github called")
		}
	},
}

func init() {
	rootCmd.AddCommand(signUpCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// signUpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// signUpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	signUpCmd.Flags().BoolVarP(&github, "github", "g", false, "Authenticate using your GitHub account")
}
