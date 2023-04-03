package conf

import (
	"log"
	"os"
)

func SetLogger(logname string) (*os.File, error) {

	// Remove existing log file
	os.Remove(logname)

	// Create log file
	logfile, err := os.OpenFile(logname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	// Set output of standard logger to log file
	log.SetOutput(logfile)

	return logfile, nil

}
