package fsutils

import (
	"fmt"
	"os"
)

/*
	This file holds functions that modify the local filesystem
*/

// Mkdir will create a directory at path, parent directories will be created
// if they do not already exist.  That is, it calls os.MkdirAll(path).
func Mkdir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("error creating directory %s => %s", path, err)
		}
	}
	return nil
}
