package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	profile "github.com/karetskiiVO/ProfileCLIUtility/pkg/profile"
	yaml "gopkg.in/yaml.v3"
)

func must0(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func must[T any](val T, err error) T {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return val
}

type namedStruct struct {
	Name string
	profile.Struct
}

func parseProfilePath(profilePath string, strict bool) ([]namedStruct, error) {
	name := strings.TrimSuffix(filepath.Base(profilePath), ".yaml")

	data, err := os.ReadFile(profilePath)
	if err != nil {
		if strict {
			return []namedStruct{}, fmt.Errorf("failed to read profile %q: %w", profilePath, err)
		} else {
			fmt.Printf("Failed to read profile %q: %v\n", profilePath, err)
			return []namedStruct{}, nil
		}
	}

	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)

	var profile profile.Struct
	if err := decoder.Decode(&profile); err != nil {
		if strict {
			return []namedStruct{}, fmt.Errorf("failed to unmarshal profile %q: %w", profilePath, err)
		} else {
			fmt.Printf("%q is not a valid profile: %v\n", profilePath, err)
			return []namedStruct{}, nil
		}
	}

	return []namedStruct{
		{
			Name:   name,
			Struct: profile,
		},
	}, nil
}
