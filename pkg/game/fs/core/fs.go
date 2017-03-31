package core

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"

	"github.com/ChrisRx/dungeonfs/pkg/game"
)

var (
	Name = "DungeonFS"
)

var GameEngine game.Engine

func init() {
}

type FS struct {
	root *Directory
}

func NewFS(root *Directory) (*FS, error) {
	r := &FS{root: root}
	return r, nil
}

func (r *FS) Root() (fs.Node, error) {
	if r.root == nil {
		return nil, errors.New("Must provide game assets")
	}
	return r.root, nil
}

func (r *FS) MountAndServe(mountpoint string, readonly bool) error {
	mountOpts := []fuse.MountOption{
		fuse.FSName(Name),
		fuse.Subtype(Name),
		fuse.VolumeName(Name),
		fuse.LocalVolume(),
	}
	if readonly {
		mountOpts = append(mountOpts, fuse.ReadOnly())
	}
	conn, err := fuse.Mount(mountpoint, mountOpts...)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch := make(chan os.Signal)
	// controlling terminal close, daemon not exit
	signal.Ignore(syscall.SIGHUP)
	signal.Notify(ch,
		os.Interrupt,
		os.Kill,
		syscall.SIGALRM,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-ch
		err := fuse.Unmount(mountpoint)
		if err != nil {
			log.Fatal(err)
		}
		if err := os.Remove(mountpoint); err != nil {
			log.Fatal(err)
		}
		conn.Close()
		os.Exit(0)
		//}
	}()

	if err = fs.Serve(conn, r); err != nil {
		return err
	}

	<-conn.Ready
	if err = conn.MountError; err != nil {
		return err
	}

	return nil
}
