package cmd

import (
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

var (
	logfile *os.File
)

func init() {
	if logfile == nil {
		var err error
		logfile, err = os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic("Failed to open log file:" + err.Error())
		}

	}
	log.SetOutput(logfile)
	log.SetLevel(log.InfoLevel)
	log.Infof("Begun %v", time.Now())
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{})

}

func Close() {
	logfile.Close()
}
