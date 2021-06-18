package log

import (
	"os"
	"strings"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

type (
	Entry interface {
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

		WithField(key string, value interface{}) Entry
	}

	Logger interface {
		Entry
	}

	logger struct {
		*logrus.Logger
	}

	entry struct {
		*logrus.Entry
	}
)

func (e *entry) WithField(key string, value interface{}) Entry {
	logrusEntry := e.Entry.WithField(key, value)
	return &entry{
		logrusEntry,
	}
}

func (l *logger) WithField(key string, value interface{}) Entry {
	logrusEntry := l.Logger.WithField(key, value)
	return &entry{
		logrusEntry,
	}
}

func New(scope string) Logger {
	l := logrus.New()
	l.SetFormatter(&nested.Formatter{
		Prefix:          "[" + scope + "]",
		TimestampFormat: "15:04:05",
	})

	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "trace":
		l.ReportCaller = true
		l.SetLevel(logrus.TraceLevel)
	case "debug":
		l.ReportCaller = true
		l.SetLevel(logrus.DebugLevel)
	case "info":
		l.SetLevel(logrus.InfoLevel)
	case "warn":
		l.SetLevel(logrus.WarnLevel)
	default:
		l.SetLevel(logrus.DebugLevel)
	}
	return &logger{
		Logger: l,
	}
}
