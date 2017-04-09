package eval

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"strconv"
)

func (c *Context) Eval(expr ast.Expr) (*reflect.Value, error) {
	switch e := expr.(type) {
	case *ast.BinaryExpr:
		return c.evalBinaryExpr(e)
	case *ast.UnaryExpr:
		return c.evalUnaryExpr(e)
	case *ast.CallExpr:
		return c.evalCallExpr(e)
	case *ast.SelectorExpr:
		return c.evalSelectorExpr(e)
	case *ast.Ident:
		return c.evalIdent(e)
	case *ast.BasicLit:
		switch e.Kind {
		//case token.INT:
		//case token.FLOAT:
		//case token.IMAG:
		//case token.CHAR:
		case token.STRING:
			s, err := strconv.Unquote(e.Value)
			if err != nil {
				panic(err)
			}
			v := reflect.ValueOf(s)
			return &v, nil
		default:
			panic(fmt.Errorf("unhandled BasicLit type '%s'", e.Kind))
		}
	default:
		panic(fmt.Errorf("unhandled type: '%v'\n", reflect.TypeOf(expr)))
	}
}

func (c *Context) evalCallExpr(e *ast.CallExpr) (*reflect.Value, error) {
	v1, err := c.Eval(e.Fun)
	if err != nil {
		return nil, err
	}
	PkgLogger.Printf("evalCallExpr: %v\n", v1)
	args := make([]reflect.Value, 0)
	for _, arg := range e.Args {
		v, err := c.Eval(arg)
		if err != nil {
			return nil, err
		}
		args = append(args, *v)
	}
	if err := checkValues(v1); err != nil {
		PkgLogger.Printf("evalCallExpr: invalid return value\n")
		return nil, err
	}
	value := v1.Call(args)
	if len(value) == 0 {
		vv := reflect.ValueOf(0)
		return &vv, nil
	}
	PkgLogger.Printf("evalCallExpr: result=%+v\n", value[0])
	return &value[0], nil
}

func (c *Context) evalIdent(e ast.Expr) (*reflect.Value, error) {
	k, ok := e.(*ast.Ident)
	if !ok {
		panic(fmt.Errorf("Expected type '*ast.Ident', received '%v'", reflect.TypeOf(e)))
	}
	if val, ok := c.identifiers[k.Name]; ok {
		PkgLogger.Printf("evalIdent: %+v, identifer=%+v\n", val, e)
		v := reflect.ValueOf(val)
		return &v, nil
	}
	PkgLogger.Printf("evalIdent: %+v\n", e)
	return nil, EvalError{Type: "UnknownIdentifier", Expr: e, Msg: fmt.Sprintf("no identifier '%s'", k.Name)}
}

func (c *Context) evalSelectorExpr(e *ast.SelectorExpr) (*reflect.Value, error) {
	val, err := c.Eval(e.X)
	if err != nil {
		return nil, err
	}
	PkgLogger.Printf("evalSelectorExpr: (%+v), valid=%t, select=%s\n", e, val.IsValid(), e.Sel.Name)
	if !val.IsValid() {
		PkgLogger.Printf("evalSelectorExpr: INVALID %+v\n", e)
		return nil, ErrInvalidValue
	}
	PkgLogger.Printf("evalSelectorExpr: (%+v), type=%v, val=%+v, select=%s\n", e, val.Type(), val, e.Sel.Name)
	PkgLogger.Printf("evalSelectorExpr: valid %+v\n", e)
	v := val.MethodByName(e.Sel.Name)
	if !v.IsValid() {
		PkgLogger.Printf("evalSelectorExpr: method=%s not valid\n", e.Sel.Name)
		for i := 0; i < val.NumField(); i++ {
			m := val.Type().Field(i)
			PkgLogger.Printf("Field %v: %+v\n", val, m.Name)
		}
		for i := 0; i < val.NumMethod(); i++ {
			m := val.Type().Method(i)
			PkgLogger.Printf("Method %v: %+v\n", val, m.Name)
		}
		return nil, ErrInvalidValue
	}
	return &v, nil
}

func (c *Context) evalBinaryExpr(e *ast.BinaryExpr) (*reflect.Value, error) {
	v1, err := c.Eval(e.X)
	if err != nil {
		return nil, err
	}
	switch e.Op {
	case token.EQL:
		v2, err := c.Eval(e.Y)
		if err != nil {
			return nil, err
		}
		r, err := compare(v1, v2)
		if err != nil {
			return nil, err
		}
		PkgLogger.Printf("compare: %t\n", r.Bool())
		return r, nil
	case token.LAND:
		if !v1.IsValid() {
			PkgLogger.Printf("isn't valid moving on %+v\n", e.X)
			vv := reflect.ValueOf(false)
			return &vv, nil
		}
		if !v1.Bool() {
			return v1, nil
		}
		v2, err := c.Eval(e.Y)
		if err != nil {
			return nil, err
		}
		return v2, nil
	case token.LOR:
		if !v1.IsValid() {
			PkgLogger.Printf("isn't valid moving on %+v\n", e.X)
			vv := reflect.ValueOf(false)
			return &vv, nil
		}
		if v1.Bool() {
			return v1, nil
		}
		v2, err := c.Eval(e.Y)
		if err != nil {
			return nil, err
		}
		return v2, nil
	default:
		panic(fmt.Errorf("Op is '%v'\n", reflect.TypeOf(e.Op)))
	}
}
func (c *Context) evalUnaryExpr(e *ast.UnaryExpr) (*reflect.Value, error) {
	v1, err := c.Eval(e.X)
	if err != nil {
		return nil, err
	}
	PkgLogger.Printf("evalUnaryExpr: %+v\n", v1)
	switch e.Op {
	case token.NOT:
		vv := reflect.ValueOf(!v1.Bool())
		return &vv, nil
	default:
		panic(fmt.Errorf("unknown unary expr '%v'", e.Op))
	}
}
