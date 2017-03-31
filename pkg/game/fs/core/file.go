package core

import (
	//"fmt"
	"os"
	"time"

	"bazil.org/fuse"
	fusefs "bazil.org/fuse/fs"
	"bazil.org/fuse/fuseutil"
	"golang.org/x/net/context"

	"github.com/ChrisRx/dungeonfs/pkg/game/fs"
)

type File struct {
	node
}

func NewFile(name string, path string) *File {
	node := NewNode(name, 0, path)
	return &File{node}
}

func (f *File) Content() []byte {
	return f.MetaData().GetBytes("Content")
}

// New always returns a nil interface because files cannot have
// child nodes in the same way as `fs.Directory` nodes.
func (f *File) New(t fs.NodeType, name string) fs.Node {
	return nil
}

func (f *File) Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fusefs.Handle, error) {
	PkgLogger.Printf("FileOpen: %s, %+v\n", f.Name(), req)
	//resp.Flags |= fuse.OpenDirectIO
	return f, nil
}

func (f *File) Getxattr(ctx context.Context, req *fuse.GetxattrRequest, resp *fuse.GetxattrResponse) error {
	PkgLogger.Printf("Filexattr: %+v\n", req)
	return nil
}

func (f *File) Attr(ctx context.Context, a *fuse.Attr) error {
	PkgLogger.Printf("FileAttr: %s, %+v\n", f.Name(), a)
	// TODO: handling the inode may be intrinsic to fuse pkg and therefore
	// unnecessary to handle explicitly here
	//a.Inode = 0
	a.Inode = f.inode
	a.Mode = 0755
	a.Uid = uint32(os.Getuid())
	a.Gid = uint32(os.Getgid())
	a.Size = uint64(len(f.Content()))
	a.Mtime = time.Now()
	return nil
}

func (f *File) Setattr(ctx context.Context, req *fuse.SetattrRequest, resp *fuse.SetattrResponse) error {
	PkgLogger.Printf("FileSetattr: %s, %+v\n", f.Name(), req)
	return nil
}

func (f *File) Fsync(ctx context.Context, req *fuse.FsyncRequest) error {
	PkgLogger.Printf("FileFsync: %+v\n", req)
	return nil
}

func (f *File) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	PkgLogger.Printf("FileRead: %+v\n", req)
	fuseutil.HandleRead(req, resp, f.Content())
	return nil
}

func (f *File) Write(ctx context.Context, req *fuse.WriteRequest, resp *fuse.WriteResponse) error {
	PkgLogger.Printf("FileWrite: %+v\n", req)
	resp.Size = len(req.Data)
	f.MetaData().Set("Content", req.Data)
	return nil
}

func (f *File) ReadAll(ctx context.Context) ([]byte, error) {
	PkgLogger.Printf("FileReadAll: %+v\n", f.Name())
	return f.Content(), nil
}
