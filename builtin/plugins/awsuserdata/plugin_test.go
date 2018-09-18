package awsuserdata

import (
	"fmt"
	"os/user"
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

	cases := map[string]struct {
		config      *AwsUserDataPluginConfig
		contains    []string
		notContains []string
	}{
		"simple": {
			config: &AwsUserDataPluginConfig{},
			notContains: []string{
				"nameserver",
			},
		},
		"with-nameservers": {
			config: &AwsUserDataPluginConfig{
				Nameservers: []string{"1.1.1.1", "2.2.2.2"},
			},
			contains: []string{
				"nameserver 1.1.1.1",
				"nameserver 2.2.2.2",
			},
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {

			_user, err := user.Current()

			check(err)

			k, err := ioutil.ReadFile(fmt.Sprintf("%s/.ssh/id_rsa.pub", _user.HomeDir))

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

			u := &AwsUserDataPlugin{
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

			name := ".user-data.sh"

			if !b.FileExists(name) {
				t.Fatalf("Expected file %s to exist", name)
			}

			bs, err := b.ReadFile(name)
			check(err)

			contents := string(bs)

			contains := append(v.contains, []string{
				fmt.Sprintf(`username=%s`, _user.Username),
				fmt.Sprintf(`ssh_public_key="%s"`, strings.TrimSuffix(string(k), "\n")),
			}...)

			for _, s := range contains {

				if !strings.Contains(contents, s) {
					t.Fatalf("user-data does not contain [%s]", s)
				}
			}

			for _, s := range v.notContains {

				if strings.Contains(contents, s) {
					t.Fatalf("user-data contains [%s]", s)
				}
			}

		})
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
			input: `
executable test {
	plugin aws-user-data {}
}
`,
			plugins: plugins,
			expected: &config.Config{
				Executables: map[string]*config.Executable{
					"test": {
						Plugins: map[string]interface{}{
							"aws-user-data": &AwsUserDataPluginConfig{},
						},
					},
				},
			}},
		"with-nameservers": {
			input: `
executable test {
	plugin aws-user-data {
		nameservers = ["1.1.1.1", "2.2.2.2"]
    }
}
`,
			plugins: plugins,
			expected: &config.Config{
				Executables: map[string]*config.Executable{
					"test": {
						Plugins: map[string]interface{}{
							"aws-user-data": &AwsUserDataPluginConfig{
								Nameservers: []string{"1.1.1.1", "2.2.2.2"},
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
