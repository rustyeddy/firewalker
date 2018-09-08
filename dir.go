package firewalker

// Directory represents a directory
type Dir struct {
	Basedir string
	Subdirs []string
	Files   []*File
	Err     error
}
