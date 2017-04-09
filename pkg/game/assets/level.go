package assets

import (
	"fmt"
	"go/parser"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/ChrisRx/dungeonfs/pkg/eval"
	"github.com/ChrisRx/dungeonfs/pkg/game/fs"
	"github.com/ChrisRx/dungeonfs/pkg/game/fs/core"
)

type Level struct {
	root           fs.Node
	objects, paths map[string]fs.Node
	properties     map[string]map[string]PropertyFunc
}

type PropertyFunc func() (*reflect.Value, error)

func NewLevel(root fs.Node, objects map[string]fs.Node) *Level {
	return &Level{
		root:       root,
		objects:    objects,
		paths:      make(map[string]fs.Node),
		properties: make(map[string]map[string]PropertyFunc),
	}
}

func (l *Level) Root() *core.Directory {
	return l.root.(*core.Directory)
}

func parseStringSlice(v interface{}) []string {
	ss := make([]string, 0)
	vv, ok := v.([]interface{})
	if !ok {
		panic(fmt.Errorf("expected type '[]interface{}', received %v", reflect.TypeOf(v)))
	}
	for _, name := range vv {
		name, ok := name.(string)
		if !ok {
			panic(fmt.Errorf("expected type 'string', received %v", reflect.TypeOf(v)))
		}
		ss = append(ss, name)
	}
	return ss
}

func (l *Level) visit(node fs.Node) {
	PkgLogger.Printf("NodeParent: %s\n", node.Name())
	for k, v := range node.MetaData().Iter() {
		switch strings.ToLower(k) {
		case "adjacent":
			fallthrough
		case "contains":
			for _, name := range parseStringSlice(v) {
				if n, ok := l.objects[name]; ok {
					PkgLogger.Printf("NodeChild: %s\n", n.Name())
					n.Parent(node)
					node.Children().Set(name, n)
					n.Path(filepath.Join(node.Path(), n.Name()))
					l.paths[n.Path()] = n
					l.visit(n)
				}
			}
		case "properties":
			for name, value := range parseAttrs(v) {
				if err := l.AddProperty(node, name, value); err != nil {
					panic(err)
				}
			}
		default:
		}
	}
}

func (l *Level) AddProperty(node fs.Node, attrName string, v interface{}) error {
	c, ok := v.(string)
	if !ok {
		return fmt.Errorf("Condition is %v, expected map[string]interface{}", reflect.TypeOf(v))
	}
	objects := make(map[string]interface{})
	for k, v := range l.objects {
		objects[k] = v
	}
	ctx, err := eval.NewContext(objects)
	if err != nil {
		return err
	}
	s := strings.Replace(c, "node", node.Name(), -1)
	s = strings.Replace(s, "\n", "", -1)
	t, err := parser.ParseExpr(s)
	if err != nil {
		return err
	}
	f := func() (*reflect.Value, error) {
		return ctx.Eval(t)
	}
	if _, ok := l.properties[node.Name()]; !ok {
		l.properties[node.Name()] = make(map[string]PropertyFunc)
	}
	l.properties[node.Name()][attrName] = f
	val, err := f()
	if err != nil {
		return err
	}
	return SetNodeAttr(node, attrName, val)
}

func (l *Level) GetProperties(key string) (map[string]PropertyFunc, bool) {
	if val, ok := l.properties[key]; ok {
		return val, ok
	}
	return nil, false
}

func SetNodeAttr(node fs.Node, k string, v *reflect.Value) error {
	if !v.IsValid() || v.Kind() == reflect.Ptr && v.IsNil() {
		return fmt.Errorf("reflect.Value is invalid")
	}
	switch v.Interface().(type) {
	case int:
		node.MetaData().Set(k, v.Int())
	case bool:
		node.MetaData().Set(k, v.Bool())
	case []byte:
		node.MetaData().Set(k, v.Bytes())
	case string:
		node.MetaData().Set(k, v.String())
	default:
		return fmt.Errorf("unhandled condition type '%v'", reflect.TypeOf(v))
	}
	return nil
}
