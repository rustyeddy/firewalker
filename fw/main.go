package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rustyeddy/fsutils"
	log "github.com/sirupsen/logrus"
)

type Stats struct {
	Files     int64
	Dirs      int64
	Others    int64
	TotalSize int64
}

var (
	action  = flag.String("action", "scan", "Actions to perform on dir")
	glob    = flag.Bool("glob", true, "Treat match as a glob (*.go, ..) ")
	pattern = flag.String("pattern", "", "Match this pattern regexp or glob")
	verbose = flag.Bool("verbose", false, "Print progress and other stuff")

	logout   = flag.String("output", "stdout", "Where to send the output from the logger")
	loglevel = flag.String("level", "warn", "Set the default level to warn")
	format   = flag.String("format", "json", "Output format color, text, JSON ... ")
	nocolors = flag.Bool("no-colors", false, "Output text logging without colors ")
)

func main() {
	flag.Parse()

	logerr := setupLogerr()

	walker := fsutils.NewWalker(getRootDirs(flag.Args()))
	walker.Verbose = true
	walker.Logerr = logerr

	// Start reading messages
	go walker.ReadMessages(os.Stderr)
	walker.StartWalking()
	fmt.Println(walker.String())
}

// Much of this can be in the library
func setupLogerr() (l *fsutils.Logerr) {
	l = fsutils.NewLogerr()
	if *logout != "stdout" {
		rd, err := os.Open(*logout)
		if err != nil {
			log.Fatalf("failed to open logerr %s", *logout)
		}
		l.SetOutput(rd)
	}

	switch *format {
	case "json":
		l.Formatter = &log.JSONFormatter{}
	case "text":
		fallthrough
	default:
		log.Errorf("expected format (json|text) got (%s) ", *format)
		l.Formatter = &log.TextFormatter{}
	}
	lvl, err := log.ParseLevel(*loglevel)
	if err != nil {
		lvl = log.WarnLevel
	}
	l.SetLevel(lvl)
	if *nocolors {
		//l.DisableColors = true
	}
	return l
}

// getRootDirs will default to current directory
func getRootDirs(d []string) (roots []string) {
	roots = flag.Args()
	if len(roots) == 0 {
		// We could default to this directory. or Fail
		// fmt.Fprintf(os.Stderr, "Need arguments to proceed ... ")
		// We will default to the local directory
		roots = []string{"."}
	}
	return roots
}
