package fsutils

import (
	"os"
)

// Directory represents a directory
type Directory struct {
	Node
	Subdirs []*Directory
	Files   []*os.FileInfo
	Err     error
}

// GetDirectory will return an object representing the
// Directory at pathname.
func GetDirectory(path string) (dir *Directory) {
	dir = &Directory{
		Node: Node{Path: path},
	}
	return dir
}
