package main

import (
	"os"

	"fmt"

	"github.com/activatedio/wrangle/builtin/plugins/ansiblehosts"
	"github.com/activatedio/wrangle/builtin/plugins/awsuserdata"
	"github.com/activatedio/wrangle/builtin/plugins/template"
	"github.com/activatedio/wrangle/command"
	"github.com/activatedio/wrangle/config"
	"github.com/activatedio/wrangle/context"
	"github.com/activatedio/wrangle/plugin"
)

func main() {

	f, err := os.Open("wrangle.hcl")
	defer f.Close()

	check(err)

	r := buildPluginRegistry()

	parser := &config.DefaultParser{
		PluginRegistry: &registryAdaptor{r},
	}

	c, err := parser.Parse(f)

	check(err)

	context, err := context.NewContext(c)

	check(err)

	var plugins []plugin.Plugin

	for k, _ := range context.Executable.Plugins {

		p, ok := r.Get(k)

		if !ok {
			panic("Invalid plugin: " + k)
		}

		plugins = append(plugins, p)
	}

	runner := command.NewDefaultRunner(context, plugins)

	err = runner.Run()

	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("Unexpected error:\n%s\n", err))
		os.Exit(1)
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func buildPluginRegistry() plugin.Registry {

	var r plugin.DefaultRegistry

	r = map[string]plugin.Plugin{
		"template":      &template.TemplatePlugin{},
		"aws-user-data": &awsuserdata.AwsUserDataPlugin{},
		"ansible-hosts": &ansiblehosts.AnsibleHostsPlugin{},
	}

	return r
}

type registryAdaptor struct {
	PluginRegistry plugin.Registry
}

func (self *registryAdaptor) Get(name string) (config.WithConfig, bool) {

	result, ok := self.PluginRegistry.Get(name)
	return result, ok
}
