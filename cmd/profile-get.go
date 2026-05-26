package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	fileutils "github.com/karetskiiVO/ProfileCLIUtility/internal/file-utils"
	cobra "github.com/spf13/cobra"
)

var profileGetCmdArgs struct {
	name string
}

var profileGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get details of a profile",
	Long:  `Get details of a profile`,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, _ []string) error {
		dir := profileCmdArgs.path.Path()
		profilePath := filepath.Join(dir, profileGetCmdArgs.name+".yaml")

		if !fileutils.Exists(profilePath) {
			return fmt.Errorf("profile not found: %s", profileGetCmdArgs.name)
		}

		_, err := os.ReadFile(profilePath)
		if err != nil {
			return fmt.Errorf("failed to read profile: %w", err)
		}

		parsedProfiles, err := parseProfilePath(profilePath, true)
		if err != nil {
			return fmt.Errorf("failed to parse profile: %w", err)
		}

		if len(parsedProfiles) == 0 {
			return fmt.Errorf("profile is empty: %s", profileGetCmdArgs.name)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

		profile := parsedProfiles[0]
		fmt.Fprintf(w, "Name\t%s\t\n", profile.Name)
		fmt.Fprintf(w, "User\t%s\t\n", profile.User)
		fmt.Fprintf(w, "Project\t%s\t\n", profile.Project)
		w.Flush()

		return nil
	},
	SilenceUsage: true,
}

func init() {
	profileGetCmd.Flags().StringVar(
		&profileGetCmdArgs.name,
		"name",
		"",
		"(required) Name of the profile to get details for",
	)
	profileGetCmd.MarkFlagRequired("name")

	profileCmd.AddCommand(profileGetCmd)
}
