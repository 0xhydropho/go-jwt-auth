package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func init() {
	Log = logrus.New()

	Log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "15:04:05.000",
		FullTimestamp:   true,
		ForceColors:     true,
	})

	Log.SetOutput(os.Stdout)
}
