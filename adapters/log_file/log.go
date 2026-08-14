package log_file

import (
	"log"
	"os"
)

type logger struct {
	log *log.Logger
}

func NewLogger(file *os.File) *logger {
	loggerDestination := log.New(file, "archimedes:", log.LstdFlags)

	return &logger{
		log: loggerDestination,
	}
}

func (l *logger) Log(msg string) {
	l.log.Println(msg)
}
