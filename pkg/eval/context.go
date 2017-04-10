package eval

import (
	"fmt"
	"go/token"
	"reflect"
)

var types = map[string]interface{}{
	"true":  true,
	"false": false,
}

type Scope struct {
	Objects map[string]interface{}
}

func NewScope() *Scope {
	return &Scope{
		Objects: make(map[string]interface{}),
	}
}

func (s *Scope) Insert(k string, v interface{}) {
	s.Objects[k] = v
}

func (s *Scope) Lookup(name string) interface{} {
	if val, ok := s.Objects[name]; ok {
		return val
	}
	return nil
}

type Context struct {
	fset        *token.FileSet
	identifiers map[string]interface{}
	returnSet   bool
	returnValue *reflect.Value
}

func NewContext(identifiers map[string]interface{}) (*Context, error) {
	c := &Context{
		fset:        token.NewFileSet(),
		identifiers: identifiers,
	}
	for k, v := range types {
		if _, ok := c.identifiers[k]; ok {
			fmt.Printf("Ident '%s' already exists\n", k)
			continue
		}
		c.identifiers[k] = v
	}
	return c, nil
}
