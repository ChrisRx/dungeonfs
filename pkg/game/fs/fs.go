package fs

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

type NodeType int

const (
	FileNode NodeType = iota
	DirNode
)

type Node interface {
	fs.Node

	New(NodeType, string) Node

	MetaData() NodeMetaData
	Name() string
	Path() string
	IsDir() bool
	IsFile() bool
	Content() []byte
	Children() Nodes

	Entry() fuse.Dirent
}

type NodeMetaData map[string]interface{}

func (n NodeMetaData) Get(key string) (interface{}, bool) {
	if val, ok := n[key]; ok {
		return val, ok
	}
	return nil, false
}

func (n NodeMetaData) GetString(key string) (v string) {
	if val, ok := n.Get(key); ok {
		v, ok = val.(string)
		// TODO: this should either panic or ignore
		if !ok {
			panic("No way")
		}
	}
	return
}

func (n NodeMetaData) GetBytes(key string) (v []byte) {
	if val, ok := n.Get(key); ok {
		v, ok = val.([]byte)
		if !ok {
			panic("No way")
		}
	}
	return
}

func (n NodeMetaData) Set(key string, value interface{}) {
	n[key] = value
}

type Nodes interface {
	Iter() []Node
	Directories() []Node
	Files() []Node
	Get(string) (Node, bool)
	Delete(string)
	Set(string, Node)
}
