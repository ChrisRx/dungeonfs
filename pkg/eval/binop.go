package eval

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

func compare(a, b *reflect.Value) (*reflect.Value, error) {
	if err := checkValues(a, b); err != nil {
		return nil, err
	}
	if a.Type() != a.Type() {
		panic(fmt.Errorf("Cannot compare, mismatch types %v and %v", a.Type(), b.Type()))
	}
	var v bool
	switch a.Interface().(type) {
	case []byte:
		if bytes.Compare(a.Bytes(), b.Bytes()) == 0 {
			v = true
		}
	case string:
		if strings.Compare(a.String(), b.String()) == 0 {
			v = true
		}
	default:
		panic(fmt.Errorf("compare '%+v' is '%v'\n", a, a.Type()))
	}
	vv := reflect.ValueOf(v)
	return &vv, nil
}
