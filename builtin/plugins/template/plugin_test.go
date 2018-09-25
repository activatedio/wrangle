package template

import (
	"testing"

	"os"

	"path/filepath"

	"strings"

	"reflect"

	"errors"

	"github.com/activatedio/wrangle/config"
	"github.com/activatedio/wrangle/context"
	"github.com/activatedio/wrangle/e2e"
	"github.com/activatedio/wrangle/plugin"
)

var _ (plugin.Plugin) = (*TemplatePlugin)(nil)

func TestTemplatePlugin(t *testing.T) {

	plugin.TestPlugin(t, &TemplatePlugin{}, &TemplatePluginConfig{
		Suffixes: map[string]*SuffixConfig{},
	},
		func(c interface{}) {
			c.(*TemplatePluginConfig).DataFile = "foo"
		})
}

func TestTemplatePlugin_Filter(t *testing.T) {

	cases := map[string]struct {
		config            *TemplatePluginConfig
		variables         map[string]string
		existsAndContains map[string][]string
		err               error
	}{
		"tmpl-simple": {
			config: &TemplatePluginConfig{
				DataFile: "data.yml",
			},
			existsAndContains: map[string][]string{
				"main.txt": []string{
					"a = \"a1\"",
					"b = \"b1\"",
				},
			},
		},
		"tmpl-variables": {
			config: &TemplatePluginConfig{
				DataFile: "data.yml",
			},
			variables: map[string]string{
				"c": "c1",
				"d": "d1",
			},
			existsAndContains: map[string][]string{
				"main.txt": []string{
					"a = \"a1\"",
					"b = \"b1\"",
					"c = \"c1\"",
					"d = \"d1\"",
				},
			},
		},
		"tmpl-variables-no-data-file": {
			config: &TemplatePluginConfig{},
			variables: map[string]string{
				"c": "c1",
				"d": "d1",
			},
			existsAndContains: map[string][]string{
				"main.txt": []string{
					"c = \"c1\"",
					"d = \"d1\"",
				},
			},
		},
		"tf-simple": {
			config: &TemplatePluginConfig{
				DataFile: "data.yml",
				Suffixes: map[string]*SuffixConfig{
					".tft": &SuffixConfig{"-generated.tf"},
				},
			},
			existsAndContains: map[string][]string{
				"main-generated.tf": []string{
					"a = \"a1\"",
					"b = \"b1\"",
				},
			},
		},
		"tf-all-funcs": {
			config: &TemplatePluginConfig{
				DataFile: "data.yml",
				Suffixes: map[string]*SuffixConfig{
					".tft": &SuffixConfig{"-generated.tf"},
				},
			},
			existsAndContains: map[string][]string{
				"main-generated.tf": []string{
					"a1 = \"b\",\"c\",\"d\"",
					"a2 = b,c,d",
					"b1 = \"d1\",\"d2\"",
				},
			},
		},
		"tf-error-1": {
			config: &TemplatePluginConfig{
				DataFile: "data.yml",
				Suffixes: map[string]*SuffixConfig{
					".tft": &SuffixConfig{"-generated.tf"},
				},
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

			c := &plugin.StubContext{
				GlobalContext: &context.Context{
					Variables: v.variables,
				},
			}

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
		"tf-simple": {
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
							"template": &TemplatePluginConfig{
								Suffixes: map[string]*SuffixConfig{},
							},
						},
					},
				},
			}},
		"with-data-file-and-one-suffix": {
			`
executable test {
	plugin template {
		data-file = "data.yml"
		suffix ".tft" {
			to-suffix = "-generated.tf"
		}
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
								Suffixes: map[string]*SuffixConfig{
									".tft": &SuffixConfig{
										ToSuffix: "-generated.tf",
									},
								},
							},
						},
					},
				},
			}},
		"with-data-file-and-two-suffixes": {
			`
executable test {
	plugin template {
		data-file = "data.yml"
		suffix ".tft" {
			to-suffix = "-generated.tf"
		}
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
								Suffixes: map[string]*SuffixConfig{
									".tft": &SuffixConfig{
										ToSuffix: "-generated.tf",
									},
								},
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
