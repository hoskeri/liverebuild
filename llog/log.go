package llog

import stdlog "log"

type Level string

const lWarn = Level("[warn] ")
const lDebug = Level("[debug] ")
const lFatal = Level("[fatal] ")
const lError = Level("[error] ")
const lInfo = Level("[info]")

func logLn(lvl Level, v ...interface{}) {
	stdlog.Print(append([]interface{}{string(lvl)}, v...)...)
}

func logF(lvl Level, fmt string, v ...interface{}) {
	stdlog.Printf(string(lvl)+" "+fmt, v...)
}

func Error(v ...interface{})            { logLn(lError, v...) }
func Errorf(f string, v ...interface{}) { logF(lError, f, v...) }

func Debug(v ...interface{})            { logLn(lDebug, v...) }
func Debugf(f string, v ...interface{}) { logF(lDebug, f, v...) }

func Warn(v ...interface{})            { logLn(lWarn, v...) }
func Warnf(f string, v ...interface{}) { logF(lWarn, f, v...) }

func Info(v ...interface{})            { logLn(lInfo, v...) }
func Infof(f string, v ...interface{}) { logF(lInfo, f, v...) }

func Fatal(v ...interface{})            { stdlog.Fatal(v...) }
func Fatalf(f string, v ...interface{}) { stdlog.Fatalf(f, v...) }

func init() {
	stdlog.SetFlags(0)
}
