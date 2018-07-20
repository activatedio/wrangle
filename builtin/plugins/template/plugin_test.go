package template

import (
	"testing"

	"os"

	"path/filepath"

	"strings"

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

	b := e2e.NewBinary("", filepath.Join("test-fixtures", "simple"))
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal("Couldn't get current working directory")
	}
	os.Chdir(b.Path())
	defer os.Chdir(orig)

	u := &TemplatePlugin{
		Config: &TemplatePluginConfig{
			DataFile: "data.yml",
		},
	}
	c := &plugin.StubContext{}

	err = u.Filter(c)

	if err != nil {
		t.Fatalf("Unexpected error %s", err)
	}

	if c.NextCallCount != 1 {
		t.Fatalf("Expected next call %d times", 1)
	}

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

}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
