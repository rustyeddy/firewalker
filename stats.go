package fsutils

import (
	"fmt"
	"os"
)

// Keep Some Status along the way
type Stats struct {
	Basedir   string
	Files     int
	Dirs      int
	Other     int
	LocalSize int64 // Just files in this directory
	TotalSize int64 // Size of local and subdirs
}

func (s *Stats) String() string {
	return fmt.Sprintf("%d files at %.1f GB\n ", s.Files, float64(s.TotalSize)/1e9)
}

// UpdateStats
func (s *Stats) Update(fi os.FileInfo) {
	s.Files++
	s.TotalSize += fi.Size()
}

func (s *Stats) Collect(fich <-chan os.FileInfo) {
	//
}
