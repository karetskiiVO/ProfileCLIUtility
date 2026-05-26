package cmd

import (
	patharg "github.com/karetskiiVO/ProfileCLIUtility/internal/path-arg"

	cobra "github.com/spf13/cobra"
)

var profileCmdArgs struct {
	path patharg.ExistsDir
}

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage profiles",
	Long:  "A command for managing profiles, including creating, viewing, listing, and deleting profiles.",
	Args:  cobra.ExactArgs(0),
}

func init() {
	must0(profileCmdArgs.path.Set("."))
	profileCmd.PersistentFlags().Var(
		&profileCmdArgs.path,
		"path",
		"Path to the profile container",
	)

	rootCmd.AddCommand(profileCmd)
}
