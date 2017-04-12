package assets

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/ChrisRx/dungeonfs/pkg/game/fs"
	"github.com/ChrisRx/dungeonfs/pkg/game/fs/core"
)

type ComponentType int

const (
	UnknownComponent ComponentType = iota // initial component state is invalid
	FileComponent                         // file component used to create file entities
	DirComponent                          // directory component used to create directory entities
	BaseComponent                         // only able to extend file/dir components
)

type Component struct {
	name  string
	t     ComponentType
	attrs map[string]interface{}
	base  string
	bases []*Component
}

func (c *Component) Name() string        { return c.name }
func (c *Component) Type() ComponentType { return c.t }
func (c *Component) Base() string        { return c.base }
func (c *Component) Bases() []*Component { return c.bases }

type Components map[string]*Component

func (c Components) LookupBaseType(t string) (*Component, bool) {
	if val, ok := c[t]; ok {
		return val, ok
	}
	return nil, false
}

func parseBaseType(s string) ComponentType {
	switch s {
	case "file":
		return FileComponent
	case "dir":
		return DirComponent
	case "base":
		return BaseComponent
	default:
		return UnknownComponent
	}
}

type Entity interface {
	fs.Node
}

// defaultAttrs defines attributes inherited by all instances of Entity
var defaultAttrs = map[string]interface{}{
	"hidden":    false,
	"permitted": true,
}

type Resources struct {
	components Components
	entities   map[string]Entity
	*Level
}

func New() *Resources {
	a := &Resources{
		components: make(Components),
		entities:   make(map[string]Entity),
	}
	return a
}

func LoadFromFile(folder string) (*core.Directory, error) {
	r := New()
	return r.LoadDir(folder)
}

func (r *Resources) GetObject(key string) Entity {
	if val, ok := r.entities[key]; ok {
		return val
	}
	return nil
}

func parseAttrs(a interface{}) map[string]interface{} {
	attrs := make(map[string]interface{})
	if a == nil {
		return attrs
	}
	aa, ok := a.(map[interface{}]interface{})
	if !ok {
		panic(fmt.Errorf("attrs wrong type: '%v'", reflect.TypeOf(a)))
	}
	for k, v := range aa {
		key, ok := k.(string)
		if !ok {
			panic("key is not string")
		}
		attrs[key] = v
	}
	return attrs
}

func (r *Resources) LoadFile(f string) error {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}
	rs := make(map[string]interface{})
	err = yaml.Unmarshal([]byte(data), &rs)
	if err != nil {
		return err
	}
	for k, v := range rs {
		parts := strings.SplitN(k, ":", 2)
		if len(parts) != 2 {
			panic("missing base type")
		}
		c := &Component{
			name:  parts[1],
			t:     parseBaseType(parts[0]),
			attrs: parseAttrs(v),
			base:  parts[0],
			bases: make([]*Component, 0),
		}
		if _, ok := r.components[c.Name()]; ok {
			panic(fmt.Errorf("component '%s' already exists"))
		}
		r.components[c.Name()] = c
		PkgLogger.Printf("RegisteredComponent: %#v\n", c)
	}
	return nil
}

func (r *Resources) LoadDir(folder string) (*core.Directory, error) {
	assetFiles, err := filepath.Glob(filepath.Join(folder, "*.yaml"))
	if err != nil {
		return nil, err
	}
	for _, f := range assetFiles {
		if err := r.LoadFile(f); err != nil {
			return nil, err
		}
	}
	unresolved := make(Components)
	for k, v := range r.components {
		if v.Type() == UnknownComponent {
			unresolved[k] = v
		}
	}
	for len(unresolved) > 0 {
		for k, v := range unresolved {
			val, ok := r.components.LookupBaseType(v.Base())
			if !ok {
				panic(fmt.Errorf("missing base component type '%s'", v.Base()))
			}
			c, ok := r.components[k]
			if !ok {
				panic(fmt.Errorf("unable to find component '%s'", k))
			}
			c.bases = append(c.bases, val)
			c.t = val.Type()
			delete(unresolved, k)
		}
	}
	for _, c := range r.components {
		switch c.Type() {
		case DirComponent:
			n := core.NewDirectory(c.Name(), nil)
			for k, v := range c.attrs {
				n.MetaData().Set(strings.ToLower(k), v)
			}
			for k, v := range defaultAttrs {
				n.MetaData().Set(strings.ToLower(k), v)
			}
			for _, base := range c.Bases() {
				//PkgLogger.Printf("Extending[%s]: %#v\n", n.Name(), base)
				for k, v := range base.attrs {
					if k == "doc" {
						continue
					}
					n.MetaData().Set(strings.ToLower(k), v)
				}
			}
			PkgLogger.Printf("Entity[%s]: %##v\n", n.Name(), n)
			r.entities[c.Name()] = n
		case FileComponent:
			n := core.NewFile(c.Name(), nil)
			for k, v := range c.attrs {
				n.MetaData().Set(strings.ToLower(k), v)
			}
			for k, v := range defaultAttrs {
				n.MetaData().Set(strings.ToLower(k), v)
			}
			for _, base := range c.Bases() {
				//PkgLogger.Printf("Extending[%s]: %#v\n", n.Name(), base)
				for k, v := range base.attrs {
					if k == "doc" {
						continue
					}
					n.MetaData().Set(strings.ToLower(k), v)
				}
			}
			PkgLogger.Printf("Entity[%s]: %##v\n", n.Name(), n)
			r.entities[c.Name()] = n
		default:
			panic("something very wrong")
		}
	}
	root := r.GetObject("Root")
	l := NewLevel(root, r.entities)
	r.Level = l
	l.visit(root)
	return l.Root(), nil
}
