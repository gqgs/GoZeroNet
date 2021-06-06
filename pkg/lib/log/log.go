package log

import (
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

type Logger interface {
	Fatal(args ...interface{})
	Error(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Debug(args ...interface{})
	Trace(args ...interface{})

	Fatalf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Tracef(format string, args ...interface{})
}

func New(scope string) Logger {
	logger := logrus.New()
	logger.SetFormatter(&nested.Formatter{
		TimestampFormat: "15:04:05",
		HideKeys:        true,
	})
	return logger.WithField("scope", scope)
}
