package game

import (
	"bazil.org/fuse"

	"github.com/ChrisRx/dungeonfs/pkg/game/fs"
)

type ActionType int

const (
	LookupAction ActionType = iota
	CreateAction
)

type Engine interface {
	Access(fs.Node) error

	// Actions checks input for possible actions to perform such as
	// when a player uses 'look' or 'use'
	Actions(ActionType, string, fs.Node) fs.Node

	Entities(fs.Node) ([]fuse.Dirent, error)
}
