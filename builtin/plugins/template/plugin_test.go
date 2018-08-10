package template

import (
	"testing"

	"os"

	"path/filepath"

	"strings"

	"reflect"

	"errors"

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
		config            *TemplatePluginConfig
		existsAndContains map[string][]string
		err               error
	}{
		"simple": {
			config: &TemplatePluginConfig{
				DataFile: "data.yml",
			},
			existsAndContains: map[string][]string{
				"main-generated.tf": []string{
					"a = \"a1\"",
					"b = \"b1\"",
				},
			},
		},
		"all-funcs": {
			config: &TemplatePluginConfig{
				DataFile: "data.yml",
			},
			existsAndContains: map[string][]string{
				"main-generated.tf": []string{
					"a1 = \"b\",\"c\",\"d\"",
					"a2 = b,c,d",
					"b1 = \"d1\",\"d2\"",
				},
			},
		},
		"error-1": {
			config: &TemplatePluginConfig{
				DataFile: "data.yml",
			},
			existsAndContains: map[string][]string{},
			err:               errors.New("template: main.tft:1: function \"regions\" not defined"),
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

			if err != nil && v.err == nil {
				t.Fatalf("Unexpected error %s", err)
			} else if !reflect.DeepEqual(err, v.err) {
				t.Fatalf("Expected [%+v], got [%+v]", v.err, err)
			}

			if err == nil && c.NextCallCount != 1 {
				t.Fatalf("Expected next call %d times", 1)
			}

			for fileName, contains := range v.existsAndContains {

				if !b.FileExists(fileName) {
					t.Fatalf("Expected file %s to exist", fileName)
				}

				bs, err := b.ReadFile(fileName)
				check(err)

				contents := string(bs)

				for _, s := range contains {

					if !strings.Contains(contents, s) {
						t.Fatalf("main.tf does not contain [%s]", s)
					}
				}
			}

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
executable test {
    plugin template {}
}
`,
			plugins,
			&config.Config{
				Executables: map[string]*config.Executable{
					"test": {
						Plugins: map[string]interface{}{
							"template": &TemplatePluginConfig{},
						},
					},
				},
			}},
		"with-data-file": {
			`
executable test {
	plugin template {
		data-file = "data.yml"
	}
}
`,
			plugins,
			&config.Config{
				Executables: map[string]*config.Executable{
					"test": {
						Plugins: map[string]interface{}{
							"template": &TemplatePluginConfig{
								DataFile: "data.yml",
							},
						},
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
