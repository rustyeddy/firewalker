package firewalker

import (
	"fmt"
	"os"
	"path/filepath"
)

// File represents a File
type File struct {
	Basedir string // Parent directory of this file
	os.FileInfo
	*Content
}

// Create and return a file.
func FileFromInfo(path string, entry os.FileInfo) *File {
	return &File{
		Basedir:  path,
		FileInfo: entry,
		Content:  nil,
	}
}

// Content contains a place holder for incoming untranslated
// text or outgoing translated text.
type Content struct {
	Buffer []byte
}

// Path returns the full path of the file object
func (f *File) Path() string {
	return filepath.Join(f.Basedir, f.Name())
}

// Read from disk
func (f *File) Read() (n int, err error) {
	if f.FileInfo == nil {
		return 0, fmt.Errorf("No FileInfo to read from")
	}
	//n, err = f.FileInfo.Read(f.Buffer)
	if err != nil {
		return 0, err
	}
	return n, nil
}

// Write content to disk
func (f *File) Write() (n int, err error) {
	if f.FileInfo == nil {
		return 0, fmt.Errorf("No FileInfo to read from")
	}
	//n, err = f.FileInfo.Write(f.Buffer)
	return n, err
}
