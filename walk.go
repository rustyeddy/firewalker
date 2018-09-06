package fsutils

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Walker struct {
	Roots []string
	Stats
	Verbose bool

	sync.WaitGroup
	FiChan  chan os.FileInfo
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
		Roots:  roots,
		Logerr: NewLogerr(),

		FiChan:  make(chan os.FileInfo),
		DirChan: make(chan string),
	}
	w.SetOutput(os.Stderr)
	w.SetLevel(log.WarnLevel)
	w.Formatter = &log.JSONFormatter{}
	return w
}

// WalkDir does a recursive walk down a directory, sending
// filesizes over the sizeChan channel.
func (w *Walker) WalkDir(path string) {

	//P("walking dir " + path)
	w.Debugln("Walking dir ", path)

	// Make sure our wait group is decremented before this
	// function returns
	defer func() {
		w.Done()
	}()

	// Loop each entry and create more subdir searches.  Making
	// sure the waitgroup is updated properly
	for _, entry := range Dirlist(path) {
		if entry.IsDir() {
			subdir := filepath.Join(path, entry.Name())

			if w.UseDirChan {
				go func() {
					P("  chan <- dir " + subdir)
					w.DirChan <- subdir // Do not block writting to channel
				}()

			} else {
				P("  walkDir " + subdir)
				w.Add(1)
				go w.WalkDir(subdir)
			}
		} else {
			P(" chan <- file " + entry.Name())
			w.FiChan <- entry
		}
	}
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
		P("  Waiting for the wait group ")
		w.Wait()
		P("  Closing File and Directory Channels ")
		close(w.FiChan)
		close(w.DirChan)
	}()

	// Create a ticker to update the user of progress.  Verbose
	// if true will cause the ticker to emit the scan summary
	// at that point.
	w.Tick = CreateTicker(500*time.Millisecond, w.Verbose)

	P("Starting Stats Loop")
	w.StatsLoop()
}

func (w *Walker) StatsLoop() {
	for {
		select {
		case fi, ok := <-w.FiChan:
			if !ok {
				return
			}
			//P("  read file channel " + fi.Name())
			w.Stats.Update(fi)
		case path, ok := <-w.DirChan:
			if !ok {
				return
			}
			w.Add(1)
			w.Stats.Dirs++
			go w.WalkDir(path)

		case _, ok := <-w.Tick:
			// the tick channel will be readable every 1sec (or ...)
			// it prints an update on os.Stdio for the user.  If
			// verbose is false, the tick channel is never written to.
			if !ok {
				w.Warn("The ticker is dead")
			}
			fmt.Println(w.Stats.String())
		}
		//fmt.Printf("%+v", w.Stats)
	}
}
