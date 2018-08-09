package config

import (
	"bytes"
	"io"

	"fmt"

	"errors"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

type Config struct {
	Executables map[string]*Executable `hcl:"executable"`
}

type Executable struct {
	Plugins map[string]interface{}
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
		Executables: make(map[string]*Executable),
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

	for _, ev := range tree.Node.(*ast.ObjectList).Items {

		if ev.Keys[0].Token.Text == "executable" {
			executableName := ev.Keys[1].Token.Text
			executable := &Executable{
				Plugins: make(map[string]interface{}),
			}
			err = hcl.DecodeObject(executable, ev.Val)

			for _, pv := range ev.Val.(*ast.ObjectType).List.Items {

				pluginName := pv.Keys[1].Token.Text
				p, ok := self.PluginRegistry.Get(pluginName)
				if !ok {
					return nil, errors.New(fmt.Sprintf("Unknown plugin %s\n", pluginName))
				}
				pc := p.GetConfig()
				err = hcl.DecodeObject(pc, pv.Val)

				executable.Plugins[pluginName] = pc
			}

			c.Executables[executableName] = executable
		}

	}

	return c, nil

}
