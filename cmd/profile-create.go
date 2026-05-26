package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	fileutils "github.com/karetskiiVO/ProfileCLIUtility/internal/file-utils"
	profile "github.com/karetskiiVO/ProfileCLIUtility/pkg/profile"

	cobra "github.com/spf13/cobra"

	yaml "gopkg.in/yaml.v3"
)

var profileCreateCmdArgs struct {
	name    string
	user    string
	project string
}

var profileCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new profile",
	Long:  "Create a new profile with specified name and details.",
	Args:  cobra.ExactArgs(0),
	RunE: func(_ *cobra.Command, _ []string) error {
		dir := profileCmdArgs.path.Path()
		name := profileCreateCmdArgs.name
		user := profileCreateCmdArgs.user
		project := profileCreateCmdArgs.project

		yamlPath := filepath.Join(dir, name+".yaml")

		if fileutils.Exists(yamlPath) {
			return fmt.Errorf("profile %q already exists at %q", name, yamlPath)
		}

		profile := profile.Struct{
			User:    user,
			Project: project,
		}

		yamlBytes, err := yaml.Marshal(&profile)
		if err != nil {
			return fmt.Errorf("failed to marshal profile: %w", err)
		}

		if err := os.WriteFile(yamlPath, yamlBytes, 0644); err != nil {
			return fmt.Errorf("failed to write profile file: %w", err)
		}

		return nil
	},
	SilenceUsage: true,
}

func init() {
	profileCreateCmd.Flags().StringVar(
		&profileCreateCmdArgs.name,
		"name",
		"",
		"(required) Name of the profile to create",
	)
	must0(profileCreateCmd.MarkFlagRequired("name"))

	profileCreateCmd.Flags().StringVar(
		&profileCreateCmdArgs.user,
		"user",
		"",
		"(required) User associated with the profile",
	)
	must0(profileCreateCmd.MarkFlagRequired("user"))

	profileCreateCmd.Flags().StringVar(
		&profileCreateCmdArgs.project,
		"project",
		"",
		"(required) Project associated with the profile",
	)
	must0(profileCreateCmd.MarkFlagRequired("project"))

	profileCmd.AddCommand(profileCreateCmd)
}
