// Package file implements core/log.ILogger by writing timestamped lines to
// an open file.
package file

import (
	"log"
	"os"
)

// logger is the core/log.ILogger implementation backing this package.
type logger struct {
	log *log.Logger
}

// NewLogger builds a logger that writes to file, prefixing every line with
// "archimedes:" and the standard date/time flags. Opening and closing file
// is the caller's responsibility.
func NewLogger(file *os.File) *logger {
	loggerDestination := log.New(file, "archimedes:", log.LstdFlags)

	return &logger{
		log: loggerDestination,
	}
}

// Log writes msg as a new line.
func (l *logger) Log(msg string) {
	l.log.Println(msg)
}
