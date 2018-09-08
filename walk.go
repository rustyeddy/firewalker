package firewalker

import (
	"fmt"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Walker struct {
	Roots []string
	Stats
	Verbose bool
	Filters []Filter
	Filemap map[string]*File

	sync.WaitGroup
	FiChan  chan *File
	DirChan chan string
	Tick    <-chan time.Time

	UseDirChan bool
	*Logerr
}

func P(str string) {
	fmt.Println(str)
}

func (w *Walker) String() string {
	return fmt.Sprintf("roots %v files %d size %d\n", w.Roots, w.Files, w.TotalSize)
}

// NewWalker will create a new directory walker for the given path
func NewWalker(roots []string) *Walker {
	w := &Walker{
		Roots:      roots,
		Logerr:     NewLogerr(),
		FiChan:     make(chan *File),
		UseDirChan: false,
	}
	if w.UseDirChan {
		w.DirChan = make(chan string)
	}
	w.SetOutput(os.Stderr)
	w.SetLevel(log.WarnLevel)
	w.Formatter = &log.JSONFormatter{}
	return w
}

// AddFilters will add a new Filter to the filter list the file
// is going to have to transform to.
func (w *Walker) AddFilter(f Filter) {
	w.Filters = append(w.Filters, f)
}

// Build a filemap
func (w *Walker) CreateFilemap() {
	w.Filemap = make(map[string]*File, 1000)
}

// WalkDir does a recursive walk down a directory, sending
// filesizes over the sizeChan channel.
func (w *Walker) WalkDir(path string) {

	w.Debugln("Walking dir ", path)

	// Make sure our wait group is decremented before this
	// function returns
	defer func() {
		w.Debugln("WG.Done Waiting")
		w.Done()
	}()

	// Loop each entry and create more subdir searches.  Making
	// sure the waitgroup is updated properly
	for _, entry := range Dirlist(path) {
		w.QEntry(path, entry)
	}
}

// QEntry will determine if the entry is a file or directory (or
// something else) and launch a new search Go routine.
func (w *Walker) QEntry(path string, entry os.FileInfo) {
	if entry.IsDir() {
		w.qdir(path, entry)
	} else {
		w.qfile(path, entry)
	}
}

func (w *Walker) qfile(path string, entry os.FileInfo) {
	w.FiChan <- FileFromInfo(path, entry)
}

// Q up the new directory and
func (w *Walker) qdir(path string, entry os.FileInfo) {
	dircomm := func(path string) {
		defer w.Done()
		w.WalkDir(path)
	}
	if w.UseDirChan {
		dircomm = func(path string) {
			defer w.Done()
			w.DirChan <- path // Do not block writting to channel
		}
	}

	w.Add(1)
	go dircomm(path)
}

// Walk through the filters transforming the file as needed
func (w *Walker) Filter(fi *File) (fout *File) {
	fout = fi
	for _, filter := range w.Filters {
		if fout = filter(w, fout); fout == nil {
			return nil
		}
	}
	return fout
}

func (w *Walker) StartWalking() {
	// Create the size channel to report file sizes, simply gather
	// sizes and total them (also count the number of files)
	// fiChan := make(chan os.FileInfo) // in walker

	// Create go routines for all roots provided from the
	for _, root := range w.Roots {
		if w.UseDirChan {
			go func() {
				w.DirChan <- root
			}()
		} else {
			w.Add(1)
			go w.WalkDir(root)
		}
	}

	// Wait for all walk functions to complete then close the
	// sizeChan, when everything completes we will.
	go func() {
		w.Debugln("  Waiting for the wait group ")
		w.Wait()
		w.Debugln("  Closing File and Directory Channels ")
		close(w.FiChan)
		if w.UseDirChan {
			close(w.DirChan)
		}
	}()

	// Create a ticker to update the user of progress.  Verbose
	// if true will cause the ticker to emit the scan summary
	// at that point.
	w.Tick = CreateTicker(500*time.Millisecond, w.Verbose)

	w.Debugln("Starting The Chain Loop")
	fok, dok, tok := true, false, true
	if w.UseDirChan {
		dok = true
	}

	for {
		var fi *File
		var path string

		select {
		case fi, fok = <-w.FiChan:
			if fok {
				w.Debugln("  read file channel " + fi.Name())
				w.Filter(fi)
			} else {
				w.Infoln("  File Channel Closed ")
			}

		case path, dok = <-w.DirChan:
			if dok {
				w.Stats.Dirs++
				w.WalkDir(path)
			} else {
				w.Infoln("  Dir Channel Closed")
			}

		case _, tok = <-w.Tick:
			// the tick channel will be readable every 1sec (or ...)
			// it prints an update on os.Stdio for the user.  If
			// verbose is false, the tick channel is never written to.
			if !tok {
				w.Warn("The ticker is dead")
			}
			fmt.Println(w.Stats.String())
		}

		if !(fok || dok) {
			fmt.Printf("~~ %+v ~~", w.Stats)
			return
		}
	}
}
