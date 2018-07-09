package llog

import stdlog "log"

type Logger struct {
	stdlog.Logger
}

func (l *Logger) Error(v ...interface{}) {
	stdlog.Print(v...)
}

func (l *Logger) Debugf(f string, v ...interface{}) {
	stdlog.Printf(f, v...)
}

func (l *Logger) Debugln(v ...interface{}) {
	stdlog.Print(v...)
}
