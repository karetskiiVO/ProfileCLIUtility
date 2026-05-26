package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	fileutils "github.com/karetskiiVO/ProfileCLIUtility/internal/file-utils"
	cobra "github.com/spf13/cobra"
)

var profileDeleteCmdArgs struct {
	verbose bool
	name    string
}

var profileDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a profile",
	Long:  `Delete a profile by name. This will remove the corresponding YAML file.`,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, _ []string) error {
		dir := profileCmdArgs.path.Path()
		profilePath := filepath.Join(dir, profileDeleteCmdArgs.name+".yaml")

		if !fileutils.Exists(profilePath) {
			return fmt.Errorf("profile not found: %s", profileDeleteCmdArgs.name)
		}

		if _, err := parseProfilePath(profilePath, true); err != nil {
			return fmt.Errorf("invalid profile to delete: %w", err)
		}

		if err := os.Remove(profilePath); err != nil {
			return fmt.Errorf("failed to delete profile: %w", err)
		}

		if profileDeleteCmdArgs.verbose {
			fmt.Printf("Profile %q deleted successfully.\n", profileDeleteCmdArgs.name)
		}

		return nil
	},
	SilenceUsage: true,
}

func init() {
	profileDeleteCmd.Flags().StringVar(
		&profileDeleteCmdArgs.name,
		"name",
		"",
		"(required) Name of the profile to delete (without .yaml extension)",
	)
	profileDeleteCmd.MarkFlagRequired("name")

	profileDeleteCmd.Flags().BoolVarP(
		&profileDeleteCmdArgs.verbose,
		"verbose",
		"v",
		false,
		"Enable verbose output",
	)

	profileCmd.AddCommand(profileDeleteCmd)
}
