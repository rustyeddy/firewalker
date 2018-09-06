package fsutils

import "os"

// CWD get the working directory
func PrintCWD() (dir string, err error) {
	return os.Getwd()
}

// CWD returns the current working directory
func CWD() (dir string, err error) {
	return os.Getwd()
}
