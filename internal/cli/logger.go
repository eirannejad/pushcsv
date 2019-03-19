package cli

import (
	"log"
)

type Logger struct {
	PrintDebug bool
}

func NewLogger() *Logger {
	return &Logger{}
}

func (m *Logger) Debug(args ...interface{}) {
	if m.PrintDebug {
		log.Print(args...)
	}
}

func (m *Logger) Print(args ...interface{}) {
	log.Print(args...)
}
