package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
)

var (
	Log *logrus.Logger
)

const (
	LogFile = "oscptools.log"
	DateISO = "2006-01-02 15:04:05"
)

func init() {
	Log = logrus.New()
	Log.Formatter = &logrus.TextFormatter{ForceColors: true, FullTimestamp: true,
		TimestampFormat: DateISO}
	Log.Out = os.Stdout

	fileWriter, err := os.OpenFile(LogFile, os.O_CREATE|os.O_WRONLY, 0600)
	if err == nil {
		Log.AddHook(&writer.Hook{
			Writer:    fileWriter,
			LogLevels: logrus.AllLevels,
		})
	}
}
