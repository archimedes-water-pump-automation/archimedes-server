package log_file

import (
	"log"
	"os"
)

type logger struct {
	log *log.Logger
}

func NewLogger(log_file string) *logger {
	f, err := os.OpenFile(log_file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil
	}
	loggerDestination := log.New(f, "archimedes:", log.LstdFlags)

	return &logger{
		log: loggerDestination,
	}
}

func (l *logger) Log(msg string) {
	l.log.Println(msg)
}
