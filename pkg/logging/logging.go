package logging

import (
	"fmt"
)

var Logger logger

func init() {
	if Logger == nil {
		Logger = &NullLogger{}
	}
}

func SetLogger(l logger) {
	Logger = l
}

type logger interface {
	Printf(string, ...interface{})
}

type NullLogger struct{}

func (l *NullLogger) Printf(format string, a ...interface{}) {
}

type DefaultLogger struct{}

func (l *DefaultLogger) Printf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}
