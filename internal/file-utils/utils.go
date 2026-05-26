package fileutils

import "os"

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func ExistsDir(path string) bool {
	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return stat.IsDir()
}
