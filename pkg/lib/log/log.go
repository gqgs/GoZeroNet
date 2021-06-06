package log

import (
	"os"

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

	WithField(key string, value interface{}) Logger
}

type logger struct {
	*logrus.Entry
}

func (l logger) WithField(key string, value interface{}) Logger {
	l.Entry = l.Entry.WithField(key, value)
	return l
}

func New(scope string) Logger {
	l := logrus.New()
	l.SetFormatter(&nested.Formatter{
		FieldsOrder:     []string{"scope"},
		TimestampFormat: "15:04:05",
		HideKeys:        true,
	})
	switch os.Getenv("LOG_LEVEL") {
	case "TRACE":
		l.SetLevel(logrus.TraceLevel)
	case "DEBUG":
		l.SetLevel(logrus.DebugLevel)
	case "INFO":
		l.SetLevel(logrus.InfoLevel)
	case "WARN":
		l.SetLevel(logrus.WarnLevel)
	default:
		l.SetLevel(logrus.InfoLevel)
	}
	return &logger{l.WithField("scope", scope)}
}
