package fsutils

import (
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
)

// FileExists makes it simple to get the answer.
func FileExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}

// FileNotExists returns true when the file does not exist
func FileNotExists(path string) bool {
	return !FileExists(path)
}

/*  -------------------------------------------------------------

    The FS utilities are all convenience functions that work on path
	strings.  Most of these functions simply wrap the underlying go
	packages calls, with some sane error checking and logging.

    -----------------------------------------------------------------*/

// ListFileInfo returns a list of *.os.FileInfo
func ListFileInfo(dir string) (files []os.FileInfo) {
	var err error
	if files, err := ioutil.ReadDir(dir); err == nil {
		return files
	}
	log.Errorf("ioutil.ReadDir failed dir %s err {%v}", dir, err)
	return nil
}

// IndexFileInfo returns a map of FileInfo, indexed on the filename (not path)
func IndexFileInfo(dir string) map[string]os.FileInfo {
	flist := ListFileInfo(dir)
	if flist == nil {
		log.Error("No index, no fun, bailing ... ", dir)
	}

	// Loop round the files creating the index
	index := make(map[string]os.FileInfo, len(flist))
	for _, f := range flist {
		index[f.Name()] = f
	}
	if len(index) < 1 {
		return nil
	}
	return index
}

// ListNames will return a array of names in directory
func ListNames(dir string) (names []string) {
	var infos []os.FileInfo
	if infos = ListFileInfo(dir); infos == nil {
		return nil
	}
	for _, fi := range infos {
		names = append(names, fi.Name())
	}
	return names
}
