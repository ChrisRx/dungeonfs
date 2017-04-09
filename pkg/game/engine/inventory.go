package engine

import (
	"github.com/ChrisRx/dungeonfs/pkg/game/fs"
	"github.com/ChrisRx/dungeonfs/pkg/game/fs/core"
)

var (
	InventoryName = ".inventory"
)

var defaultAttrs = map[string]interface{}{
	"hidden":      false,
	"permitted":   false,
	"Description": "Magic Bag of Storage",
}

type Item struct {
	Name    string
	Content []byte
}

type Inventory struct {
	*core.Directory
}

func NewInventory(parent fs.Node, items ...Item) *Inventory {
	d := core.NewDirectory(InventoryName, parent)
	for k, v := range defaultAttrs {
		d.MetaData().Set(k, v)
	}
	inv := &Inventory{d}
	inv.Add(items...)
	return inv
}

func (inv *Inventory) Add(items ...Item) {
	for _, item := range items {
		f := inv.New(fs.FileNode, item.Name)
		f.MetaData().Set("Content", item.Content)
	}
}

func (inv *Inventory) Remove(name string) {
	inv.Children().Delete(name)
}

func (inv *Inventory) Register(node fs.Node) {
	if _, ok := node.Children().Get(InventoryName); ok || node.Name() == InventoryName {
		return
	}
	node.Children().Set(inv.Name(), inv)
}

func (inv *Inventory) Unregister(node fs.Node) {
	if _, ok := node.Children().Get(InventoryName); !ok {
		return
	}
	node.Children().Delete(inv.Name())
}
