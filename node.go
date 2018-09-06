package fsutils

import (
	"fmt"
	"os"
)

// Node represents an object in the file system, either a directory
// or a file.
type Node struct {
	Path        string
	os.FileInfo // If this is nil, we are not synced with FS
	Err         *Message
}

// Info returns the file info for this node, that will tell us
// wether this node exists, and if so, all its stats, created,
// size, etc.
func (n *Node) Info() (fi os.FileInfo, fserr *Message) {
	var err error
	if n.FileInfo == nil {
		if n.FileInfo, err = os.Stat(n.Path); err != nil {
			n.Err = ErrorMessage(fmt.Errorf(n.Path + err.Error()))
		}
	}
	return n.FileInfo, n.Err
}
