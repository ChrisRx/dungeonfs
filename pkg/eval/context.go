package eval

import (
	"fmt"
)

var types = map[string]interface{}{
	"true":  true,
	"false": false,
}

type Context struct {
	identifiers map[string]interface{}
}

func NewContext(identifiers map[string]interface{}) (*Context, error) {
	c := &Context{
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
