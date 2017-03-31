package core

import (
	"os"
	"time"

	"bazil.org/fuse"
	fusefs "bazil.org/fuse/fs"
	"golang.org/x/net/context"

	"github.com/ChrisRx/dungeonfs/pkg/game"
	"github.com/ChrisRx/dungeonfs/pkg/game/fs"
)

var currentPath string

type Directory struct {
	node
	// game.Engine
}

func NewDirectory(name, path string) *Directory {
	node := NewNode(name, os.ModeDir, path)
	d := &Directory{
		node: node,
	}
	return d
}

func (d *Directory) New(t fs.NodeType, name string) fs.Node {
	switch t {
	case fs.FileNode:
		return d.NewFile(name)
	case fs.DirNode:
		return d.NewDirectory(name)
	default:
		panic("idk")
	}
}

func (d *Directory) NewDirectory(name string) *Directory {
	newDir := NewDirectory(name, d.Path())
	d.Set(name, newDir)
	return newDir
}

func (d *Directory) NewFile(name string) *File {
	newFile := NewFile(name, d.Path())
	d.Set(name, newFile)
	return newFile
}

func (d *Directory) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = d.inode
	a.Uid = uint32(os.Getuid())
	a.Gid = uint32(os.Getgid())
	a.Mode = os.ModeDir | 0755
	return nil
}

func (d *Directory) Access(ctx context.Context, req *fuse.AccessRequest) error {
	PkgLogger.Printf("DirAccess: name=%s, path=%s, req=%+v\n", d.Name(), d.Path(), req)
	return GameEngine.Access(d)
}

func (d *Directory) Lookup(ctx context.Context, req *fuse.LookupRequest, resp *fuse.LookupResponse) (fusefs.Node, error) {
	PkgLogger.Printf("DirLookup: name=%s, path=%s\n", d.Name(), d.Path())
	PkgLogger.Printf("DirLookup: %+v\n\t%+v\n", req, resp)
	resp.EntryValid = 0 * time.Second
	if action := GameEngine.Actions(game.LookupAction, req.Name, d); action != nil {
		return action, nil
	}
	if n, ok := d.Get(req.Name); ok {
		return n, nil
	}
	return nil, fuse.ENOENT
}

func (d *Directory) Create(ctx context.Context, req *fuse.CreateRequest, resp *fuse.CreateResponse) (fusefs.Node, fusefs.Handle, error) {
	PkgLogger.Printf("DirCreate: name=%s, path=%s\n", d.Name(), d.Path())
	f := GameEngine.Actions(game.CreateAction, req.Name, d)
	return f, f, nil
}

func (d *Directory) Remove(ctx context.Context, req *fuse.RemoveRequest) error {
	PkgLogger.Printf("DirRemove: name=%s, path=%s\n", d.Name(), d.Path())
	if f, ok := d.Get(req.Name); ok && f.IsFile() {
		d.Delete(req.Name)
	}
	return nil
}

func (d *Directory) Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fusefs.Handle, error) {
	PkgLogger.Printf("DirOpen: name=%s, path=%s\n", d.Name(), d.Path())
	// This isn't important because the engine keeping track of the player isn't necessary, rather the location is implicit to the requests being made. This below is a hard way of doing this and presents a situation where the actions cannot be disambiguated (i.e. between 'ls' or tab complete etc)
	//if d.Path != GameEngine.Player().CurrentPath {
	var updated bool
	if d.Path() != currentPath {
		updated = true
		currentPath = d.Path()
	}
	PkgLogger.Printf("DirOpen: currentpath=%s, updated=%t\n", currentPath, updated)
	return d, nil
}

func (d *Directory) Mkdir(ctx context.Context, req *fuse.MkdirRequest) (fusefs.Node, error) {
	PkgLogger.Printf("Mkdir: %+v\n", req.Name)
	if dir, ok := d.Get(req.Name); ok && dir.IsDir() {
		return nil, fuse.ENOENT
	}
	newDir := d.NewDirectory(req.Name)
	d.Set(req.Name, newDir)
	return newDir, nil
}

func (d *Directory) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	PkgLogger.Printf("ReadDirAll: name=%s, path=%s\n", d.Name(), d.Path())
	return GameEngine.Entities(d)
}

func (d *Directory) Flush(ctx context.Context, req *fuse.FlushRequest) error {
	PkgLogger.Printf("Flush: %+v\n", req)
	return nil
}

func (d *Directory) ReadAll(ctx context.Context) ([]byte, error) {
	PkgLogger.Printf("ReadAll: %+v\n", d.Name())
	return []byte{}, nil
}

func (d *Directory) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	PkgLogger.Printf("Read: %+v\n", d.Name())
	return nil
}
