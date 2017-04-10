package eval

import (
	"fmt"
	"go/ast"
	"reflect"
)

func (c *Context) evalDecl(decl ast.Decl) (*reflect.Value, error) {
	switch d := decl.(type) {
	case *ast.FuncDecl:
		return c.evalFuncDecl(d)
	case *ast.GenDecl:
		return c.evalGenDecl(d)
	default:
		panic(fmt.Errorf("unhandled ast.Decl type: '%v'\n", reflect.TypeOf(decl)))
	}
}

func (c *Context) evalFuncDecl(decl ast.Decl) (*reflect.Value, error) {
	return &False, nil
}

func (c *Context) evalGenDecl(decl ast.Decl) (*reflect.Value, error) {
	return &False, nil
}
