package core

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"bazil.org/fuse"

	"github.com/ChrisRx/dungeonfs/pkg/game/fs"
)

type node struct {
	inode    uint64
	name     string
	mode     os.FileMode
	path     string
	metadata nodeMetaData

	parent fs.Node
	nodes
}

func NewNode(name string, mode os.FileMode, parent fs.Node) node {
	n := node{
		inode:    NewInode(),
		name:     name,
		mode:     mode,
		parent:   parent,
		nodes:    make(nodes),
		metadata: make(nodeMetaData),
	}
	if n.parent != nil {
		n.path = filepath.Join(n.parent.Path(), name)
	} else {
		n.path = name
	}
	return n
}

func (n *node) Name() string {
	return n.name
}

func (n *node) Path(args ...string) string {
	if len(args) > 0 {
		n.path = args[0]
	}
	return n.path
}

func (n *node) Content() []byte {
	return n.MetaData().GetBytes("Content")
}

func (n *node) Parent(args ...fs.Node) fs.Node {
	if len(args) > 0 {
		n.parent = args[0]
	}
	return n.parent
}

func (n *node) Delete() bool {
	if n.parent == nil {
		return false
	}
	n.parent.Children().Delete(n.Name())
	return true
}

func (n *node) IsDir() bool {
	return n.mode&os.ModeDir != 0
}

func (n *node) IsFile() bool {
	return n.mode&os.ModeDir == 0
}

func (n *node) Children() fs.Nodes {
	return n.nodes
}

func (n *node) MetaData() fs.NodeMetaData {
	return n.metadata
}

func (n *node) Entry() fuse.Dirent {
	ent := fuse.Dirent{
		Name: n.Name(),
	}
	if n.IsDir() {
		ent.Type = fuse.DT_Dir
	} else {
		ent.Type = fuse.DT_File
	}
	return ent
}

var inode uint64

func NewInode() uint64 {
	inode += 1
	return inode
}

type nodes map[string]fs.Node

func (n nodes) Iter() []fs.Node {
	l := make([]fs.Node, 0)
	for _, v := range n {
		l = append(l, v)
	}
	return l
}

func (n nodes) Directories() []fs.Node {
	d := make([]fs.Node, 0)
	for _, v := range n.Iter() {
		if v.IsDir() {
			d = append(d, v)
		}
	}
	return d
}

func (n nodes) Exists(key string) bool {
	_, ok := n.Get(key)
	fmt.Printf("Exists '%s': %t\n", key, ok)
	fmt.Printf("Exists: '%+v\n", n)
	return ok
}

func (n nodes) Files() []fs.Node {
	f := make([]fs.Node, 0)
	for _, v := range n.Iter() {
		if v.IsDir() {
			continue
		}
		f = append(f, v)
	}
	return f
}

func (n nodes) Get(key string) (fs.Node, bool) {
	if val, ok := n[key]; ok {
		return val, ok
	}
	return nil, false
}

func (n nodes) Delete(key string) {
	delete(n, key)
}

func (n nodes) Set(key string, node fs.Node) {
	n[key] = node
}

type nodeMetaData map[string]interface{}

func (n nodeMetaData) Iter() map[string]interface{} {
	return n
}

func (n nodeMetaData) Get(key string) (interface{}, bool) {
	key = strings.ToLower(key)
	if val, ok := n[key]; ok {
		return val, ok
	}
	return nil, false
}

func (n nodeMetaData) GetString(key string) (v string) {
	if val, ok := n.Get(key); ok {
		switch val.(type) {
		case string:
			v = val.(string)
		case []byte:
			v = string(val.([]byte))
		default:
			panic(fmt.Errorf("expected type 'string' but received type '%s'", reflect.TypeOf(val)))
		}
	}
	return
}

func (n nodeMetaData) GetBool(key string) bool {
	if val, ok := n.Get(key); ok {
		vv := reflect.ValueOf(val)
		switch v := vv.Interface().(type) {
		case int:
			return v == 0
		case bool:
			return v
		}
	}
	return false
}

func (n nodeMetaData) GetBytes(key string) []byte {
	return []byte(n.GetString(key))
}

func (n nodeMetaData) Set(key string, value interface{}) {
	key = strings.ToLower(key)
	n[key] = value
}
