package engine

import (
	"fmt"
	"os"
	"time"

	"bazil.org/fuse"
	"bazil.org/fuse/fuseutil"
	"golang.org/x/net/context"
)

//type item interface {
//Name() string
//Look() string
//Use(...string) string
//}

type item struct {
	name        string
	description string
}

func NewItem(name, description string) *item {
	return &item{name, description}
}

func (i *item) Name() string {
	return i.name
}
func (i *item) Look() string {
	return i.description
}
func (i *item) Use(a ...string) string {
	return i.description
}

type items struct {
	items []*item
}

func (i *items) Count() int {
	return len(i.items)
}

func (i *items) Get(name string) (*item, bool) {
	for _, i := range i.items {
		if i.Name() == name {
			return i, true
		}
	}
	return nil, false
}

func (i *items) List() string {
	var s string
	for _, i := range i.items {
		s += fmt.Sprintf("%s\n", i.Name())
	}
	return s
}

type Inventory struct {
	items
	//*fs.File
}

func NewInventory(newItems ...*item) *Inventory {
	return &Inventory{items{newItems}}
}

func (inv *Inventory) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 2
	a.Mode = 0755
	a.Uid = uint32(os.Getuid())
	a.Gid = uint32(os.Getgid())
	var size uint64
	for _, item := range inv.items.items {
		size += uint64(len(item.Name()))
	}
	a.Size = size
	a.Mtime = time.Now()
	return nil
}

func (inv *Inventory) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	fuseutil.HandleRead(req, resp, []byte(inv.items.List()))
	return nil
}

func (inv *Inventory) ReadAll(ctx context.Context) ([]byte, error) {
	return []byte(inv.items.List()), nil
}

func (inv *Inventory) Content() []byte {
	return []byte(inv.items.List())
}

func (inv *Inventory) Name() string {
	return "inventory"
}
