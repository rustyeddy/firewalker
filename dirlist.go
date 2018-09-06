package fsutils

import (
	"io/ioutil"
	"os"
	"sync"
)

type ReadErrors sync.Map

// Add a read error for later printing
func (re ReadErrors) Add(path string, err error) {
	//re[path] = err
}

var (
	MaxConcurrentDirs int // max con
	CountingSemaphore chan struct{}
	readErrors        ReadErrors
)

// Create our semaphore based on number of concurrent dirs
func init() {
	// XXX - Must set after flags are parsed in main() after
	// this init function has been run.
	MaxConcurrentDirs = 20
}

// GetSortedEntries takes a path string and returns three arrays
// (lists) of: regular files, directories and other (pipes, perm
// denied, etc.)
func GetSortedDirlist(path string) (files, dirs, other []os.FileInfo) {
	entries := Dirlist(path)
	if entries == nil {
		return nil, nil, nil
	}
	files, dirs, other = SortDirlist(entries)
	return files, dirs, other
}

// GetEntries converts a path string to an []os.FileInfo each FileInfo
// represents one the the "paths" children.  They can be files,
// [sub]Directories or "other" dependending on the respective file
// type.  Directories may be used for deeper search (or not), files
// may be used to information gathering, translation, copy, move or
// delte, etc.
func Dirlist(path string) (entries []os.FileInfo) {
	var err error
	if CountingSemaphore == nil {
		CountingSemaphore = make(chan struct{}, MaxConcurrentDirs)
	}

	// Limit the number of dir lists we can have running at a time
	CountingSemaphore <- struct{}{}
	defer func() { <-CountingSemaphore }()

	entries, err = ioutil.ReadDir(path)
	if err != nil {
		// readErrors.Add(path, err)
		return nil
	}
	return entries
}

// SortDirlist walks through the entire directory list, identifies
// each entry as a: file, dir or other grouping them in the
// appropriate list.
func SortDirlist(dirlist []os.FileInfo) (files, dirs, other []os.FileInfo) {
	for _, fi := range dirlist {
		mode := fi.Mode()
		switch {
		case mode.IsDir():
			dirs = append(dirs, fi) // Could send to a new channel?
		case mode.IsRegular():
			files = append(files, fi)
		default:
			// TODO ~ Categorize the "other" category with a map?
			other = append(other, fi)
		}
	}
	return files, dirs, other
}
