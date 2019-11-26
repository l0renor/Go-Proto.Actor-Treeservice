package logger

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type Logger struct {
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
}

// nolint:gochecknoglobals
var instance *Logger

// nolint:gochecknoglobals
var once sync.Once

func GetInstance() *Logger {
	once.Do(func() {
		instance = &Logger{}
		Init(instance, ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	})
	return instance
}

func Init(lo *Logger, traceHandle io.Writer, infoHandle io.Writer, warningHandle io.Writer, errorHandle io.Writer) {
	lo.Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	lo.Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	lo.Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	lo.Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
