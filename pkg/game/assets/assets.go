package assets

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/ChrisRx/dungeonfs/pkg/game/fs"
	"github.com/ChrisRx/dungeonfs/pkg/game/fs/core"
)

type ResourceType int

const (
	FileResource ResourceType = iota
	DirResource
)

type Resource struct {
	name  string
	t     ResourceType
	attrs map[string]interface{}
}

func (r *Resource) Name() string       { return r.name }
func (r *Resource) Type() ResourceType { return r.t }

func parseBaseType(s string) (ResourceType, error) {
	switch s {
	case "file":
		return FileResource, nil
	case "dir":
		return DirResource, nil
	}
	return 0, fmt.Errorf("unable to parse base type")
}

type Resources struct {
	resources map[string]*Resource
	objects   map[string]fs.Node
	*Level
}

func New() *Resources {
	a := &Resources{
		resources: make(map[string]*Resource),
		objects:   make(map[string]fs.Node),
	}
	return a
}

func LoadFromFile(folder string) (*core.Directory, error) {
	r := New()
	return r.LoadDir(folder)
}

func (r *Resources) GetObject(key string) fs.Node {
	if val, ok := r.objects[key]; ok {
		return val
	}
	return nil
}

func parseAttrs(a interface{}) map[string]interface{} {
	aa, ok := a.(map[interface{}]interface{})
	if !ok {
		panic("attrs wrong type")
	}
	attrs := make(map[string]interface{})
	for k, v := range aa {
		key, ok := k.(string)
		if !ok {
			panic("key is not string")
		}
		attrs[key] = v
	}
	return attrs
}

func (a *Resources) LoadFile(f string) ([]*Resource, error) {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	rs := make(map[string]interface{})
	err = yaml.Unmarshal([]byte(data), &rs)
	if err != nil {
		return nil, err
	}
	rr := make([]*Resource, 0)

	for k, v := range rs {
		parts := strings.SplitN(k, ":", 2)
		if len(parts) != 2 {
			panic("missing base type")
		}
		name := parts[1]
		rt, err := parseBaseType(parts[0])
		if err != nil {
			panic(err)
		}
		attrs := parseAttrs(v)
		r := &Resource{
			name:  name,
			t:     rt,
			attrs: attrs,
		}
		rr = append(rr, r)
		PkgLogger.Printf("r: %+v\n", r)
	}
	return rr, nil
}

var defaultAttrs = map[string]interface{}{
	"hidden":    false,
	"permitted": true,
}

func (r *Resources) LoadDir(folder string) (*core.Directory, error) {
	assetFiles, err := filepath.Glob(filepath.Join(folder, "*.yaml"))
	if err != nil {
		return nil, err
	}
	for _, f := range assetFiles {
		rr, err := r.LoadFile(f)
		if err != nil {
			return nil, err
		}
		for _, res := range rr {
			r.resources[res.Name()] = res
		}
	}
	for _, v := range r.resources {
		switch v.Type() {
		case DirResource:
			n := core.NewDirectory(v.Name(), nil)
			for k, v := range v.attrs {
				PkgLogger.Printf("Metadata[%s]: %v: %v\n", n.Name(), k, v)
				n.MetaData().Set(strings.ToLower(k), v)
			}
			for k, v := range defaultAttrs {
				n.MetaData().Set(strings.ToLower(k), v)
			}
			r.objects[v.Name()] = n
		case FileResource:
			n := core.NewFile(v.Name(), nil)
			for k, v := range v.attrs {
				PkgLogger.Printf("Metadata[%s]: %v: %v\n", n.Name(), k, v)
				n.MetaData().Set(strings.ToLower(k), v)
			}
			for k, v := range defaultAttrs {
				n.MetaData().Set(strings.ToLower(k), v)
			}
			r.objects[v.Name()] = n
		default:
			// we might not have the type yet
		}
	}
	root := r.GetObject("Root")
	l := NewLevel(root, r.objects)
	r.Level = l
	l.visit(root)
	return l.Root(), nil
}
