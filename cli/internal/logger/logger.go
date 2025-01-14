package logger

import (
	"log"
	"os"
)

// InitLogs Initialises the logger
func InitLogs(logdir string) {

	LogFile := logdir + "/gofish.log"
	logFile, err := os.OpenFile(LogFile, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		os.Exit(1)
	}

	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile | log.LstdFlags)

}
