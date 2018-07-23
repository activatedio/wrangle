package awsuserdata

import (
	"fmt"
	user2 "os/user"
	"testing"

	"strings"

	"os"
	"path/filepath"

	"io/ioutil"

	"github.com/activatedio/wrangle/config"
	"github.com/activatedio/wrangle/e2e"
	"github.com/activatedio/wrangle/plugin"
)

var _ (plugin.Plugin) = (*AwsUserDataPlugin)(nil)

func TestTemplatePlugin_Filter(t *testing.T) {

	user, err := user2.Current()

	check(err)

	k, err := ioutil.ReadFile(fmt.Sprintf("%s/.ssh/id_rsa.pub", user.HomeDir))

	if err != nil {
		t.Skip("SSH Key not availble for test")
	}

	b := e2e.NewBinary("", filepath.Join("test-fixtures", "simple"))
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal("Couldn't get current working directory")
	}
	os.Chdir(b.Path())
	defer os.Chdir(orig)

	u := &AwsUserDataPlugin{}
	c := &plugin.StubContext{}

	err = u.Filter(c)

	if err != nil {
		t.Fatalf("Unexpected error %s", err)
	}

	if c.NextCallCount != 1 {
		t.Fatalf("Expected next call %d times", 1)
	}

	name := ".user-data.sh"

	if !b.FileExists(name) {
		t.Fatalf("Expected file %s to exist", name)
	}

	bs, err := b.ReadFile(name)
	check(err)

	contents := string(bs)

	contains := []string{
		fmt.Sprintf(`username=%s`, user.Username),
		fmt.Sprintf(`ssh_public_key="%s"`, strings.TrimSuffix(string(k), "\n")),
	}

	for _, s := range contains {

		if !strings.Contains(contents, s) {
			t.Fatalf("main.tf does not contain [%s]", s)
		}
	}

}
func TestAwsUserDataPlugin_GetConfig(t *testing.T) {

	plugins := map[string]config.WithConfig{
		"aws-user-data": &AwsUserDataPlugin{},
	}

	cases := map[string]struct {
		input    string
		plugins  map[string]config.WithConfig
		expected *config.Config
	}{
		"simple": {
			`
plugin aws-user-data {}
`,
			plugins,
			&config.Config{
				Plugins: map[string]interface{}{
					"aws-user-data": &AwsUserDataPluginConfig{},
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
