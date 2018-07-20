package config_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/activatedio/wrangle/config"
	"github.com/activatedio/wrangle/plugin"
	"github.com/davecgh/go-spew/spew"
)

var _ config.Parser = (*config.DefaultParser)(nil)

type PluginConfigA struct {
	A     string
	B     string
	Child map[string]*PluginChildConfigA
}

type PluginConfigB struct {
	E string
	F string
}

type PluginChildConfigA struct {
	C string
	D string
}

type StubRegistry map[string]config.WithConfig

func (self StubRegistry) Get(name string) (config.WithConfig, bool) {
	r, ok := self[name]
	return r, ok
}

func TestDefaultParser_Parse(t *testing.T) {

	cases := map[string]struct {
		input    string
		expexted interface{}
		registry map[string]config.WithConfig
	}{
		"empty": {"", &config.Config{
			Plugins: make(map[string]interface{}),
		},
			map[string]config.WithConfig{},
		},
		"standard": {`
delegate = "test-delegate"
plugin a {
    a = "a"
    b = "b"
	child a {
		c = "c1"
		d = "d1"
	}
	child b {
		c = "c2"
		d = "d2"
	}
}
plugin b {
    e = "e"
    f = "f"
}
`,
			&config.Config{
				Delegate: "test-delegate",
				Plugins: map[string]interface{}{
					"a": &PluginConfigA{
						A: "a",
						B: "b",
						Child: map[string]*PluginChildConfigA{
							"a": {
								C: "c1",
								D: "d1",
							},
							"b": {
								C: "c2",
								D: "d2",
							},
						},
					},
					"b": &PluginConfigB{
						E: "e",
						F: "f",
					},
				},
			},
			map[string]config.WithConfig{
				"a": &plugin.StubPlugin{
					Config: &PluginConfigA{},
				},
				"b": &plugin.StubPlugin{
					Config: &PluginConfigB{},
				},
			},
		},
	}

	for k, v := range cases {

		t.Run(k, func(t *testing.T) {

			var r StubRegistry
			r = v.registry

			u := &config.DefaultParser{
				PluginRegistry: r,
			}

			got, err := u.Parse(strings.NewReader(v.input))

			if err != nil {
				t.Fatalf("Unexpected error %s", err)
			}

			if !reflect.DeepEqual(got, v.expexted) {
				t.Fatal(spew.Sprintf("Wanted %+v, got %+v", v.expexted, got))
			}

		})

	}

}
