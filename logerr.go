package fsutils

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

// Logerr will log error and debug message in a way that
// will have little effect on the performance of the program
// overall.
type Logerr struct {
	*log.Logger
	Logch chan Message
	Errch chan Message
}

// Get a new logger
func NewLogerr() (l *Logerr) {
	l = &Logerr{
		Logger: log.New(),
	}
	l.SetOutput(os.Stderr)
	l.SetLevel(log.DebugLevel)
	l.Formatter = &log.JSONFormatter{}

	l.Logch = make(chan Message, 5)
	l.Errch = make(chan Message, 4)

	return l
}

// ReadMessages will read the Log and Error channels picking up
// and printing/logging messages as they arrive.
func (l *Logerr) ReadMessages(wr io.Writer) {
	for {
		l.SetOutput(wr)
		select {
		case msg, ok := <-l.Logch:
			if !ok {
				l.Error("Log channel has closed")
			}
			l.Print(msg)
		case msg, ok := <-l.Errch:
			if !ok {
				l.Error("Error channel has closed")
			}
			l.Error(msg)
		}
	}
}

// ---------------------- Message --------------------------

// Message is used to pass through the log and error channels
type Message struct {
	Path     string // Optional
	Filename string // Optional file the error occured in
	Lineno   int    // Optional the line the error occured on
	Message  string // Mandatory something
	Level    log.Level
	Err      error
}

// Return an error is there one
func (m *Message) Error() string {
	if m.Err == nil {
		return ""
	}
	return m.Err.Error()
}

// Return if this is an error message or not
func (m *Message) IsError() bool {
	return (m.Level >= log.ErrorLevel)
}

func ErrorMessage(err error) (m *Message) {
	m = &Message{Err: err, Level: log.ErrorLevel}
	return m
}

func WarnMessage(err error) (m *Message) {
	m = &Message{Err: err, Level: log.WarnLevel}
	return m
}

func (m *Message) String() string {
	return m.Message
}
