package logging

import (
	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Logger
)

// logrus level
// -- Trace
// -- Debug
// -- Info
// -- Warn
// -- Error
// -- Fatal
// -- Panic

func init() {
	log = logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
	})

	log.SetLevel(logrus.DebugLevel)
	log.SetReportCaller(false)
}

func GetLogger(level string) func(string, ...interface{}) {
	if level == logrus.ErrorLevel.String() {
		return log.Errorf
	}

	return log.Tracef
}

func SetLevel(env string) {
	var level string
	if env == "prod" {
		level = "info"
	} else {
		level = "debug"
	}
	if level != "debug" && level != "info" && level != "warn" && level != "error" && level != "fatal" && level != "panic" {
		log.SetLevel(logrus.DebugLevel)
		log.Warnf("Invalid log level: %s. Set to default level: debug", level)
		return
	}

	l, err := logrus.ParseLevel(level)
	if err != nil {
		log.SetLevel(logrus.DebugLevel)
	} else {
		log.SetLevel(l)
	}
}

func Tracef(format string, args ...interface{}) {
	log.Tracef(format, args...)
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	log.Panicf(format, args...)
}
