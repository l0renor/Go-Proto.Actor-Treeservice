package logger

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type logger struct {
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
}

var instance *logger
var once sync.Once

func GetInstance() *logger {
	once.Do(func() {
		instance = &logger{}
		Init(instance, ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	})
	return instance
}

func Init(
	lo *logger,
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

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
