package engine

import (
	"bytes"

	"bazil.org/fuse"

	"github.com/ChrisRx/dungeonfs/pkg/game"
	"github.com/ChrisRx/dungeonfs/pkg/game/assets"
	"github.com/ChrisRx/dungeonfs/pkg/game/fs"
)

type Engine struct {
	// game.Player
	// game.State
}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) Access(node fs.Node) error {
	if node.Name() == "door" {
		if f, ok := node.Children().Get("lock"); ok && f.IsFile() {
			if bytes.Compare(f.Content(), assets.Key) == 0 {
				node.Children().Delete("lock")
				return nil
			}
			return fuse.EPERM
		}
	}
	return nil
}

func (e *Engine) Actions(t game.ActionType, name string, node fs.Node) fs.Node {
	switch t {
	case game.LookupAction:
		return lookupAction(name, node)
	case game.CreateAction:
		return createAction(name, node)
	}
	return nil
}

func (e *Engine) Entities(node fs.Node) ([]fuse.Dirent, error) {
	l := make([]fuse.Dirent, 0)
	// Have to check access here as well in cases where ReadDirAll is called for
	// tab complete
	if err := e.Access(node); err != nil {
		for _, f := range node.Children().Files() {
			l = append(l, f.Entry())
		}
		return l, nil
	}
	// TODO: Inventory will implement fs.Node and have convenient methods to access. will have a global state
	// added to the new global state container so it will be consistent through the game
	if _, ok := node.Children().Get(".inventory"); !ok {
		newDir := node.New(fs.DirNode, ".inventory")
		newDir.MetaData().Set("Description", "An adventure's bag o' goods")
		key := newDir.New(fs.FileNode, "key")
		key.MetaData().Set("Content", assets.Key)
		newDir.New(fs.FileNode, "bean")
		newDir.New(fs.FileNode, "sword")
	}
	for _, v := range node.Children().Iter() {
		l = append(l, v.Entry())
	}
	return l, nil
}
