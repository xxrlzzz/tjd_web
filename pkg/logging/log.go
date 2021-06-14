package logging

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"traffic_jam_direction/pkg/file"
)

type Level int

var (
	F *os.File

	DefaultPrefix      = ""
	DefaultCallerDepth = 2

	logger     *log.Logger
	logPrefix  = ""
	levelFlags = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
)

// Setup initialize the log instance
func Setup(isDev bool) {
	var err error
	filePath := getLogFilePath()
	fileName := getLogFileName()
	F, err = file.MustOpen(fileName, filePath, !isDev)
	if err != nil {
		log.Fatalf("logging.Setup err: %v", err)
	}

	logger = log.New(F, DefaultPrefix, log.LstdFlags)
}

// Debug output logs at debug level
func Debug(v ...interface{}) {
	if logger == nil {
		return
	}
	setPrefix(DEBUG)
	logger.Println(v...)
}

// Info output logs at info level
func Info(v ...interface{}) {
	if logger == nil {
		return
	}
	setPrefix(INFO)
	logger.Println(v...)
}

// Warn output logs at warn level
func Warn(v ...interface{}) {
	if logger == nil {
		return
	}
	setPrefix(WARNING)
	logger.Println(v...)
}

// Error output logs at error level
func Error(v ...interface{}) {
	if logger == nil {
		return
	}
	setPrefix(ERROR)
	logger.Println(v...)
}

// Fatal output logs at fatal level
func Fatal(v ...interface{}) {
	if logger == nil {
		return
	}
	setPrefix(FATAL)
	logger.Fatalln(v...)
}

// DebugF output logs at debug level
func DebugF(fmt string, v ...interface{}) {
	if logger == nil {
		return
	}
	setPrefix(DEBUG)
	logger.Printf(fmt, v...)
}

// InfoF output logs at info level
func InfoF(fmt string, v ...interface{}) {
	if logger == nil {
		return
	}
	setPrefix(INFO)
	logger.Printf(fmt, v...)
}

// WarnF output logs at warn level
func WarnF(fmt string, v ...interface{}) {
	if logger == nil {
		return
	}
	setPrefix(WARNING)
	logger.Printf(fmt, v...)
}

// Error output logs at error level
func ErrorF(fmt string, v ...interface{}) {
	if logger == nil {
		return
	}
	setPrefix(ERROR)
	logger.Printf(fmt, v...)
}

// Fatal output logs at fatal level
func FatalF(fmt string, v ...interface{}) {
	if logger == nil {
		return
	}
	setPrefix(FATAL)
	logger.Printf(fmt, v...)
}

// setPrefix set the prefix of the log output
func setPrefix(level Level) {
	_, file, line, ok := runtime.Caller(DefaultCallerDepth)
	if ok {
		logPrefix = fmt.Sprintf("[%s][%s:%d]", levelFlags[level], filepath.Base(file), line)
	} else {
		logPrefix = fmt.Sprintf("[%s]", levelFlags[level])
	}

	logger.SetPrefix(logPrefix)
}
