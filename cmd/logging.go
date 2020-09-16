package cmd

import (
	"io/ioutil"
	"log"
	"os"
)

func init() {
	initConfig()
}

type Logging struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	debugLog      *log.Logger
}

func NewLogging(enabled bool) *Logging {

	var infoLog *log.Logger
	var debugLog *log.Logger
	var errorLog *log.Logger

	if enabled {
		infoLog = log.New(openLogFile("pgcli_info.log"), "INFO\t", log.Ldate|log.Ltime)
		debugLog = log.New(openLogFile("pgcli_debug.log"), "DEBUG\t", log.Ldate|log.Ltime)
		errorLog = log.New(openLogFile("pgcli_error.log"), "ERROR\t", log.Ldate|log.Ltime)
	} else {
		infoLog = log.New(ioutil.Discard, "", 0)
		errorLog = log.New(ioutil.Discard, "", 0)
		debugLog = log.New(ioutil.Discard, "", 0)
	}


	app := &Logging{
		errorLog: errorLog,
		infoLog: infoLog,
		debugLog: debugLog,
	}

	return app
}

func openLogFile(path string) *os.File {
	f, err := os.OpenFile(path, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	return f
}
