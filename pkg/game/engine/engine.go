package engine

import (
	//"reflect"

	"bazil.org/fuse"

	"github.com/ChrisRx/dungeonfs/pkg/game"
	"github.com/ChrisRx/dungeonfs/pkg/game/assets"
	"github.com/ChrisRx/dungeonfs/pkg/game/fs"
)

type Engine struct {
	*assets.Level
	*Player
}

func NewEngine(r *assets.Resources) *Engine {
	p := NewPlayer()
	key := r.GetObject("key")
	p.Inventory.Add(
		Item{Name: "key", Content: key.Content()},
		Item{Name: "bean"},
		Item{Name: "sword"},
	)
	return &Engine{
		Player: p,
		Level:  r.Level,
	}
}

func (e *Engine) computeProperties(node fs.Node) {
	if val, ok := e.GetProperties(node.Name()); ok {
		for k, fn := range val {
			v, err := fn()
			if err != nil {
				panic(err)
			}
			assets.SetNodeAttr(node, k, v)
		}
	}
}

func (e *Engine) Access(node fs.Node) error {
	e.computeProperties(node)
	e.Player.Register(node)
	if !node.MetaData().GetBool("permitted") {
		return fuse.EPERM
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
	for _, n := range node.Children().Iter() {
		e.computeProperties(n)
		if n.MetaData().GetBool("hidden") {
			continue
		}
		l = append(l, n.Entry())
	}
	return l, nil
}
