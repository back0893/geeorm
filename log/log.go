package log

import (
	"io"
	"log"
	"os"
	"sync"
)

var (
	textred   = "\033[0;31m"
	textgreen = "\033[0;32m"
	textreset = "\033[0m"
	errorLog  = log.New(os.Stderr, textred+"[error]"+textreset, log.LstdFlags|log.Lshortfile)
	infoLog   = log.New(os.Stderr, textgreen+"[info]"+textreset, log.LstdFlags|log.Lshortfile)
	loggers   = []*log.Logger{errorLog, infoLog}
	mu        = &sync.Mutex{}
)

var (
	Error  = errorLog.Print
	Info   = infoLog.Print
	Errorf = errorLog.Printf
	Infof  = infoLog.Printf
)

const (
	LevelInfo = iota
	LevelError
	Disabled
)

func SetLogLevel(level int) {
	mu.Lock()
	defer mu.Unlock()
	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)
	}
	if level > LevelInfo {
		infoLog.SetOutput(io.Discard)
	}
	if level > LevelError {
		errorLog.SetOutput(io.Discard)
	}
}
