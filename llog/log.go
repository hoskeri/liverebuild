package llog

import stdlog "log"

type Level string

const lDbug = Level("[debug] ")
const lInfo = Level("[ info] ")

var Verbose = false

func logLn(lvl Level, v ...interface{}) {
	if lvl == lDbug && !Verbose {
		return
	}
	stdlog.Print(append([]interface{}{string(lvl)}, v...)...)
}

func logF(lvl Level, fmt string, v ...interface{}) {
	if lvl == lDbug && !Verbose {
		return
	}
	stdlog.Printf(string(lvl)+" "+fmt, v...)
}

func Info(v ...interface{})            { logLn(lInfo, v...) }
func Infof(f string, v ...interface{}) { logF(lInfo, f, v...) }

func Debug(v ...interface{})            { logLn(lDbug, v...) }
func Debugf(f string, v ...interface{}) { logF(lDbug, f, v...) }

func Fatal(v ...interface{})            { stdlog.Fatal(v...) }
func Fatalf(f string, v ...interface{}) { stdlog.Fatalf(f, v...) }

func init() {
	stdlog.SetFlags(0)
}
