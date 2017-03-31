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

//type DefaultLogger struct{}

//func (l *DefaultLogger) Printf(format string, a ...interface{}) {
//fmt.Printf(format, a...)
//}

type DefaultLogger struct {
	last     string
	repeated int
}

func (l *DefaultLogger) Printf(format string, a ...interface{}) {
	cur := fmt.Sprintf(format, a...)
	if cur == l.last {
		l.repeated++
		fmt.Printf("\rPrevious message repeated %d times", l.repeated)
		return
	}
	if l.repeated > 0 {
		l.repeated = 0
		fmt.Printf("\n")
	}
	print(cur)
	l.last = cur
}
