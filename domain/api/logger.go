package api

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

type LogLevel uint
type Color string

const (
	Reset  Color = "\033[0m"
	Prefix       = "[API] "
	Flags        = log.LstdFlags
)

const (
	errorLevel LogLevel = iota + 1
	warnLevel
	infoLevel
	debugLevel
)

var (
	LoggerLevel = NewKey("LOGGER_LEVEL", "gte=0,lte=4", "4")

	level = debugLevel
	err   = log.New(os.Stderr, Prefix, Flags)
	out   = log.New(os.Stdout, Prefix, Flags)
)

type LogMessage struct {
	Tag   string
	Color Color
	Level LogLevel
}

var (
	errorMessage = &LogMessage{"ERROR: ", "\033[31m", errorLevel}
	warnMessage  = &LogMessage{"WARN: ", "\033[33m", warnLevel}
	infoMessage  = &LogMessage{"INFO: ", "\033[32m", infoLevel}
	debugMessage = &LogMessage{"DEBUG: ", "\033[37m", debugLevel}
)

func internalPrint(l *log.Logger, info *LogMessage, v ...interface{}) {
	if info.Level <= level {
		line := ""
		if level == debugLevel {
			if _, file, no, ok := runtime.Caller(2); ok {
				line = fmt.Sprintf("%s:%d", file, no)
			}
		}
		l.Println(line, info.Color, info.Tag, v, Reset)
	}
}

func D(v ...interface{}) {
	internalPrint(out, debugMessage, v...)
}

func I(v ...interface{}) {
	internalPrint(out, infoMessage, v...)
}

func W(v ...interface{}) {
	internalPrint(out, warnMessage, v...)
}

func E(v ...interface{}) {
	internalPrint(err, errorMessage, v...)
}

func Fatal(v ...interface{}) {
	internalPrint(err, errorMessage, v...)
	os.Exit(1)
}

func InitLogger(configuration Configuration) {
	loggerLevel := configuration.Uint(LoggerLevel)
	level = LogLevel(loggerLevel)
}
