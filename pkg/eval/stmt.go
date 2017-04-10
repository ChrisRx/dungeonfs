package eval

import (
	"fmt"
	"go/ast"
	"reflect"
)

func (c *Context) evalStmt(stmt ast.Stmt) (*reflect.Value, error) {
	if c.returnSet {
		c.returnSet = false
		if !c.returnValue.IsValid() {
			PkgLogger.Printf("WTF!: %+v\n", stmt)
		}
		return c.returnValue, nil
	}
	switch s := stmt.(type) {
	case *ast.BlockStmt:
		return c.evalBlockStmt(s)
	case *ast.IfStmt:
		return c.evalIfStmt(s)
	case *ast.ReturnStmt:
		return c.evalReturnStmt(s)
	case *ast.AssignStmt:
		return c.evalAssignStmt(s)
	case *ast.ExprStmt:
		return c.evalExprStmt(s)
	default:
		panic(fmt.Errorf("unhandled ast.Stmt type: '%v'\n", reflect.TypeOf(stmt)))
	}
}

func (c *Context) evalExprStmt(e *ast.ExprStmt) (*reflect.Value, error) {
	v, err := c.Eval(e.X)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (c *Context) evalAssignStmt(e *ast.AssignStmt) (*reflect.Value, error) {
	rhs, err := c.Eval(e.Rhs[0])
	if err != nil {
		return nil, err
	}
	name := e.Lhs[0].(*ast.Ident).Name
	c.identifiers[name] = rhs
	return &True, nil
}

func (c *Context) evalReturnStmt(e *ast.ReturnStmt) (*reflect.Value, error) {
	if len(e.Results) == 0 {
		return nil, EvalError{Type: "ReturnError", Node: e, Msg: fmt.Sprintf("not enough return arguments")}
	}
	PkgLogger.Printf("evalReturnStmt: %+v\n", e.Results[0])
	result, err := c.Eval(e.Results[0])
	if err != nil {
		return nil, err
	}
	c.returnSet = true
	c.returnValue = result
	return result, nil
}

func (c *Context) evalIfStmt(e *ast.IfStmt) (*reflect.Value, error) {
	cond, err := c.Eval(e.Cond)
	if err != nil {
		return nil, err
	}
	if cond.Bool() {
		v, err := c.Eval(e.Body)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
	return &False, nil
}

func (c *Context) evalBlockStmt(e *ast.BlockStmt) (*reflect.Value, error) {
	PkgLogger.Printf("evalBlockStmt: %+v\n", e)
	for i, stmt := range e.List {
		v, err := c.Eval(stmt)
		if err != nil {
			return nil, err
		}
		PkgLogger.Printf("Block[%d]: %+v\n", i, v)
		if c.returnSet {
			c.returnValue = v
			return v, nil
		}
	}
	return &Nil, nil
}
