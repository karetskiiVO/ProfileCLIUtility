package cmd

import (
	"fmt"
	"os"

	cobra "github.com/spf13/cobra"
)

const (
	programName = "profile-utility"
)

var rootCmd = &cobra.Command{
	Use:   programName,
	Short: "Profile utility",
	Long: `A CLI for working with profiles that supports
creating, viewing, listing, and deleting profiles,
as well as outputting help for available commands.`,
	Args: cobra.ExactArgs(0),
	SilenceErrors: true,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
