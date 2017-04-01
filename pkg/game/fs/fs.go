package fs

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

type NodeType int

const (
	FileNode NodeType = iota
	DirNode
	TempFileNode
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

type NodeMetaData interface {
	Get(string) (interface{}, bool)
	GetString(string) string
	GetBytes(string) []byte
	Set(string, interface{})
}

type Nodes interface {
	Iter() []Node
	Directories() []Node
	Files() []Node
	Get(string) (Node, bool)
	Delete(string)
	Set(string, Node)
}
