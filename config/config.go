package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

type Config struct {
	Delegate string
	Plugins  map[string]interface{}
}

type Parser interface {
	Parse(io.Reader) (*Config, error)
}

type Registry interface {
	Get(name string) (WithConfig, bool)
}

type WithConfig interface {
	GetConfig() interface{}
}

type DefaultParser struct {
	PluginRegistry Registry
}

func (self *DefaultParser) Parse(r io.Reader) (*Config, error) {

	c := &Config{
		Plugins: make(map[string]interface{}),
	}

	b := new(bytes.Buffer)
	b.ReadFrom(r)
	s := b.String()

	tree, err := hcl.Parse(s)

	if err != nil {
		return nil, err
	}

	err = hcl.DecodeObject(c, tree)

	if err != nil {
		return nil, err
	}

	for _, v := range tree.Node.(*ast.ObjectList).Items {

		if v.Keys[0].Token.Text == "plugin" {
			name := v.Keys[1].Token.Text

			p, ok := self.PluginRegistry.Get(name)

			if !ok {
				return nil, errors.New(fmt.Sprintf("Unknown plugin %s\n", name))
			}

			pc := p.GetConfig()
			err = hcl.DecodeObject(pc, v.Val)

			c.Plugins[name] = pc

		}

	}

	return c, nil

}
