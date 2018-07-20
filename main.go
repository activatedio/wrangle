package main

import (
	"os"

	"github.com/activatedio/wrangle/builtin/plugins/template"
	"github.com/activatedio/wrangle/config"
	"github.com/activatedio/wrangle/plugin"
)

func main() {

	f, err := os.Open("wrangle.hcl")
	defer f.Close()

	// We just die hard here
	check(err)
	parser := &config.DefaultParser{}

	_, err = parser.Parse(f)
	//_ := buildPluginRegistry()

	// TODO - Handle order of plugins
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func buildPluginRegistry() plugin.Registry {

	var r plugin.DefaultRegistry

	r = map[string]plugin.Plugin{
		"template": &template.TemplatePlugin{},
	}

	return r
}
