package eval

import (
	"errors"
	"fmt"
	"go/ast"
	"reflect"
)

// TODO: change to const enum
var (
	ErrInvalidValue = errors.New("reflect.Value is invalid")
)

type EvalError struct {
	ast.Expr
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
