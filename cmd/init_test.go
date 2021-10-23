package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	logfile *os.File
)

func init() {
	if logfile == nil {
		var err error
		logfile, err = ioutil.TempFile("", "example")
		if err != nil {
			panic("Failed to open log file:" + err.Error())
		}
		fmt.Println(logfile.Name())
	}
	log.SetOutput(logfile)
	log.SetLevel(log.InfoLevel)
	log.Infof("Begun %v %s", time.Now(), logfile.Name())
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{})

}

func testLogClose() {
	logfile.Close()
}
