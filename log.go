package main

import (
	"io"
	"log"
)

const (
	logFlag = log.LstdFlags | log.Lshortfile | log.Lmsgprefix | log.Lmicroseconds
)

type Logger struct {
	stdout *log.Logger
	stderr *log.Logger
}

func NewLogger(appName string, stdout, stderr io.Writer) Logger {
	l := Logger{
		stdout: log.New(stdout, appName, logFlag),
		stderr: log.New(stderr, appName, logFlag),
	}
	return l
}

func (l *Logger) Err(err error) {
	l.stderr.Output(1, " [ERR ] "+err.Error())
}
