package core

import (
	"os"
	"path/filepath"

	"bazil.org/fuse"

	"github.com/ChrisRx/dungeonfs/pkg/game/fs"
)

type node struct {
	inode    uint64
	name     string
	mode     os.FileMode
	path     string
	metadata nodeMetaData

	nodes
}

func NewNode(name string, mode os.FileMode, path string) node {
	return node{
		inode:    NewInode(),
		name:     name,
		mode:     mode,
		path:     filepath.Join(path, name),
		nodes:    make(nodes),
		metadata: make(nodeMetaData),
	}
}

func (n *node) Name() string {
	return n.name
}

func (n *node) Path() string {
	return n.path
}

func (n *node) Content() []byte {
	return n.MetaData().GetBytes("Content")
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

func (n nodeMetaData) Get(key string) (interface{}, bool) {
	if val, ok := n[key]; ok {
		return val, ok
	}
	return nil, false
}

func (n nodeMetaData) GetString(key string) (v string) {
	if val, ok := n.Get(key); ok {
		v, ok = val.(string)
		// TODO: this should either panic or ignore
		if !ok {
			panic("No way")
		}
	}
	return
}

func (n nodeMetaData) GetBytes(key string) (v []byte) {
	if val, ok := n.Get(key); ok {
		v, ok = val.([]byte)
		if !ok {
			panic("No way")
		}
	}
	return
}

func (n nodeMetaData) Set(key string, value interface{}) {
	n[key] = value
}
