package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/tabwriter"

	"github.com/karetskiiVO/ProfileCLIUtility/internal/synctool"
	cobra "github.com/spf13/cobra"
)

var profileListCmdArgs struct {
	strict   bool
	verbose  bool
	jobCount uint
}

var profileListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all profiles",
	Long:  `List all profiles in the container, just look for "*.yaml" files.`,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, _ []string) error {
		dir := profileCmdArgs.path.Path()
		pattern := filepath.Join(dir, "*.yaml")
		profilePaths := must(filepath.Glob(pattern))

		validNamedOptProfiles, err := synctool.MapPoolSlice(
			cmd.Context(),
			4,
			profilePaths,
			parseProfilePathExternalStrict,
		)

		if err != nil {
			return fmt.Errorf("failed to parse profiles: %w", err)
		}

		validNamedProfiles := slices.Concat(validNamedOptProfiles...)

		slices.SortFunc(validNamedProfiles, func(a, b namedStruct) int {
			return strings.Compare(a.Name, b.Name)
		})

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "Name\tUser\tProject\t")
		fmt.Fprintln(w, "----\t----\t-------\t")

		for _, np := range validNamedProfiles {
			fmt.Fprintf(w, "%s\t%s\t%s\t\n", np.Name, np.User, np.Project)
		}

		w.Flush()
		fmt.Println("----")
		fmt.Printf("Total: %d\n", len(validNamedProfiles))

		return nil
	},
	SilenceUsage: true,
}

func init() {
	profileListCmd.Flags().BoolVar(
		&profileListCmdArgs.strict,
		"strict",
		false,
		"If set, check that all profiles are valid and report errors for invalid profiles",
	)
	profileListCmd.Flags().BoolVarP(
		&profileListCmdArgs.verbose,
		"verbose",
		"v",
		false,
		"Set verbosity mode",
	)
	profileListCmd.Flags().UintVar(
		&profileListCmdArgs.jobCount,
		"jobs",
		1,
		"Number of concurrent jobs to use when parsing profiles",
	)

	profileCmd.AddCommand(profileListCmd)
}

func parseProfilePathExternalStrict(profilePath string) ([]namedStruct, error) {
	return parseProfilePath(profilePath, profileListCmdArgs.strict)
}
