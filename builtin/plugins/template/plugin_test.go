package template

import (
	"testing"

	"os"

	"path/filepath"

	"strings"

	"github.com/activatedio/wrangle/config"
	"github.com/activatedio/wrangle/e2e"
	"github.com/activatedio/wrangle/plugin"
)

var _ (plugin.Plugin) = (*TemplatePlugin)(nil)

func TestTemplatePlugin(t *testing.T) {

	plugin.TestPlugin(t, &TemplatePlugin{}, &TemplatePluginConfig{},
		func(c interface{}) {
			c.(*TemplatePluginConfig).DataFile = "foo"
		})
}

func TestTemplatePlugin_Filter(t *testing.T) {

	cases := map[string]struct {
		config *TemplatePluginConfig
		verify func(t *testing.T, b *e2e.Binary)
	}{
		"simple": {
			config: &TemplatePluginConfig{
				DataFile: "data.yml",
			},
			verify: func(t *testing.T, b *e2e.Binary) {
				name := "main.tft"

				if !b.FileExists(name) {
					t.Fatalf("Expected file %s to exist", name)
				}

				bs, err := b.ReadFile("main.tf")
				check(err)

				contents := string(bs)

				contains := []string{
					"a = \"a1\"",
					"b = \"b1\"",
				}

				for _, s := range contains {

					if !strings.Contains(contents, s) {
						t.Fatalf("main.tf does not contain [%s]", s)
					}
				}

			},
		},
		"all-funcs": {
			config: &TemplatePluginConfig{
				DataFile: "data.yml",
			},
			verify: func(t *testing.T, b *e2e.Binary) {
				name := "main.tft"

				if !b.FileExists(name) {
					t.Fatalf("Expected file %s to exist", name)
				}

				bs, err := b.ReadFile("main.tf")
				check(err)

				contents := string(bs)

				contains := []string{
					"a1 = \"b\",\"c\",\"d\"",
				}

				for _, s := range contains {

					if !strings.Contains(contents, s) {
						t.Fatalf("main.tf does not contain [%s]", s)
					}
				}

			},
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			b := e2e.NewBinary("", filepath.Join("test-fixtures", k))
			orig, err := os.Getwd()
			if err != nil {
				t.Fatal("Couldn't get current working directory")
			}
			os.Chdir(b.Path())
			defer os.Chdir(orig)

			u := &TemplatePlugin{
				Config: v.config,
			}

			c := &plugin.StubContext{}

			err = u.Filter(c)

			if err != nil {
				t.Fatalf("Unexpected error %s", err)
			}

			if c.NextCallCount != 1 {
				t.Fatalf("Expected next call %d times", 1)
			}

			v.verify(t, b)

		})
	}

}

func TestTemplatePlugin_Config(t *testing.T) {

	plugins := map[string]config.WithConfig{
		"template": &TemplatePlugin{},
	}

	cases := map[string]struct {
		input    string
		plugins  map[string]config.WithConfig
		expected *config.Config
	}{
		"simple": {
			`
plugin template {}
`,
			plugins,
			&config.Config{
				Plugins: map[string]interface{}{
					"template": &TemplatePluginConfig{},
				},
			}},
		"with-data-file": {
			`
plugin template {
	data-file = "data.yml"
}
`,
			plugins,
			&config.Config{
				Plugins: map[string]interface{}{
					"template": &TemplatePluginConfig{
						DataFile: "data.yml",
					},
				},
			}},
	}

	for k, v := range cases {

		t.Run(k, func(t *testing.T) {
			config.TestConfig(t, strings.NewReader(v.input), v.plugins, v.expected)
		})
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
