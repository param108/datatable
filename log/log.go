package log


import(
	"os"
	log "github.com/sirupsen/logrus"
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
			panic("Failed to open log file:" +err.Error())
		}
	}
	log.SetOutput(logfile)

	log.Infof("Begun %v", time.Now())

}

func Close() {
	logfile.Close()
}
