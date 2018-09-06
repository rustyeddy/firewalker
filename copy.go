package fsutils

import (
	"fmt"
	"io"
	"os"
)

// Copy does just that, copies the file from src to dst
func Copy(src, dst string) (err error) {

	sfi, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sfi.Mode().IsRegular() {
		// Can not copy non-regular files
		return fmt.Errorf("can not copy non-regular file: %s", src)
	}

	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("Failed to open source file %s - %v", src, err)
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create dst %s - %v", dst, err)
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}() // parens are necessary to EXECUTE teh function (deffered)

	// Now for the actual copy from one place to another
	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	// Make sure we are all sync'd
	err = out.Sync()
	return err
}
