package main

import (
	"flag"
	"fmt"
	"os"

	fw "github.com/rustyeddy/firewalker"
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
	format   = flag.String("format", "text", "Output format color, text, JSON ... ")
	nocolors = flag.Bool("no-colors", false, "Output text logging without colors ")
)

func main() {
	flag.Parse()

	walker := fw.NewWalker(getRootDirs(flag.Args()))
	walker.Verbose = false
	walker.Logerr = SetupLogerr()

	walker.AddFilter(fw.IndexOnce)
	walker.AddFilter(fw.CollectStats)

	// Start reading messages
	go walker.ReadMessages(os.Stderr)

	walker.Infof("Start walking roots %+v...\n", walker.Roots)
	walker.StartWalking()

	walker.Infoln("Main is existing, Goodbye .!. ")
	fmt.Println("\t" + walker.String())
}

// getRootDirs will default to current directory, unless
func getRootDirs(d []string) (roots []string) {
	roots = []string{"."}
	if len(d) > 0 {
		roots = d
	}
	return roots
}
