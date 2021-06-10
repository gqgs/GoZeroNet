package log

import (
	"os"
	"strings"

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

	IfError(args ...interface{})

	WithField(key string, value interface{}) Logger
}

type logger struct {
	*logrus.Entry
}

func (l logger) WithField(key string, value interface{}) Logger {
	l.Entry = l.Entry.WithField(key, value)
	return l
}

func (l logger) IfError(args ...interface{}) {
	for _, arg := range args {
		if arg != nil {
			l.Error(arg)
		}
	}
}

func New(scope string) Logger {
	l := logrus.New()
	l.SetFormatter(&nested.Formatter{
		FieldsOrder:     []string{"scope"},
		TimestampFormat: "15:04:05",
	})
	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "trace":
		l.SetLevel(logrus.TraceLevel)
	case "debug":
		l.SetLevel(logrus.DebugLevel)
	case "info":
		l.SetLevel(logrus.InfoLevel)
	case "warn":
		l.SetLevel(logrus.WarnLevel)
	default:
		l.SetLevel(logrus.DebugLevel)
	}
	return &logger{l.WithField("scope", scope)}
}
