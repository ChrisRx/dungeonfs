package eval

import (
	"errors"
	"fmt"
	"go/ast"
	"reflect"
)

var (
	False = reflect.ValueOf(false)
	True  = reflect.ValueOf(true)
	Nil   = reflect.ValueOf(nil)

	ErrInvalidValue = errors.New("reflect.Value is invalid")
	ErrReturnValue  = errors.New("return value is set")
)

func (c *Context) Eval(node ast.Node) (*reflect.Value, error) {
	switch n := node.(type) {
	case ast.Decl:
		return c.evalDecl(n)
	case ast.Expr:
		return c.evalExpr(n)
	case ast.Stmt:
		return c.evalStmt(n)
	default:
		panic(fmt.Errorf("unhandled ast.Node type: '%v'\n", reflect.TypeOf(node)))
	}
}

type EvalError struct {
	ast.Node
	Type string
	Msg  string
}

func (e EvalError) Error() string {
	return fmt.Sprintf("%s[%d:%d]: %s", e.Type, e.Pos(), e.End(), e.Msg)
}

func checkValues(values ...*reflect.Value) error {
	for _, v := range values {
		if !v.IsValid() || v.Kind() == reflect.Ptr && v.IsNil() {
			return fmt.Errorf("reflect.Value is invalid")
		}
	}
	return nil
}
