package core

type logger interface {
	Printf(string, ...interface{})
}

type NullLogger struct{}

func (l *NullLogger) Printf(format string, a ...interface{}) {
}

var PkgLogger logger

func init() {
	if PkgLogger == nil {
		PkgLogger = &NullLogger{}
	}
}
