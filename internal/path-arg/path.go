package patharg

import (
	"fmt"
	"os"
	"path/filepath"

	pflag "github.com/spf13/pflag"
)

type ExistsPath struct {
	path string
}

var _ pflag.Value = (*ExistsPath)(nil)

func (ep ExistsPath) String() string {
	return ep.path
}

func (ep *ExistsPath) Set(value string) error {
	_, err := os.Stat(value)
	if err != nil {
		return fmt.Errorf("convert %q to path: %w", value, err)
	}
	*ep = ExistsPath{path: value}
	return nil
}

func (ExistsPath) Type() string {
	return "exists_path"
}

func (ep ExistsPath) Path() string {
	res, _ := filepath.Abs(ep.path)
	return res
}

type ExistsDir struct {
	path string
}

var _ pflag.Value = (*ExistsDir)(nil)

func (ed ExistsDir) String() string {
	return ed.path
}

func (ed *ExistsDir) Set(value string) error {
	stat, err := os.Stat(value)
	if os.IsNotExist(err) {
		return fmt.Errorf("convert %q to dir: %w", value, err)
	}
	if !stat.IsDir() {
		return fmt.Errorf("convert %q to dir: not a directory", value)
	}
	*ed = ExistsDir{path: value}
	return nil
}

func (ExistsDir) Type() string {
	return "exists_dir"
}

func (ed ExistsDir) Path() string {
	res, _ := filepath.Abs(ed.path)
	return res
}
