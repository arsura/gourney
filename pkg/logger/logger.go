package logger

import (
	"log"
	"os"
)

type Logger struct {
	Info  *log.Logger
	Error *log.Logger
}

func InitLog() *Logger {
	return &Logger{
		Info:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		Error: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
}
