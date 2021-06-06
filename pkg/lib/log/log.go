package log

import "github.com/sirupsen/logrus"

type Logger interface {
	Fatal(args ...interface{})
	Error(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Debug(args ...interface{})
	Trace(args ...interface{})
}

func New(scope string) Logger {
	return logrus.New().WithField("scope", scope)
}
