package engine

import (
	"github.com/ChrisRx/dungeonfs/pkg/logging"
)

type logger interface {
	Printf(string, ...interface{})
}

var PkgLogger logger

func init() {
	if PkgLogger == nil {
		PkgLogger = &logging.NullLogger{}
	}
}
