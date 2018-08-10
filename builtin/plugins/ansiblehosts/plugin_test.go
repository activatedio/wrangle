package ansiblehosts

import (
	"strings"
	"testing"

	"os"
	"path/filepath"
	"reflect"

	"os/exec"

	"io/ioutil"

	"github.com/activatedio/wrangle/config"
	"github.com/activatedio/wrangle/context"
	"github.com/activatedio/wrangle/e2e"
	"github.com/activatedio/wrangle/plugin"
)

var _ (plugin.Plugin) = (*AnsibleHostsPlugin)(nil)

func TestTemplatePlugin_Filter(t *testing.T) {

	cases := map[string]struct {
		config *AnsibleHostsPluginConfig
		err    error
	}{
		"empty-output": {
			config: &AnsibleHostsPluginConfig{
				Modules: []string{"a"},
			},
		},
		"two-instances": {
			config: &AnsibleHostsPluginConfig{
				Modules:           []string{"a"},
				FqdnOutputName:    "instance_fqdns",
				RecordsOutputName: "instance_records",
			},
		},
		"instances-with-groups": {
			config: &AnsibleHostsPluginConfig{
				Modules:           []string{"a"},
				FqdnOutputName:    "instance_fqdns",
				RecordsOutputName: "instance_records",
				Groups: map[string]*Group{
					"build": {"^build"},
					"ops":   {"^ops"},
					"zero":  {"0"},
					"one":   {"1"},
				},
			},
		},
		"instances-with-groups-two-modules": {
			config: &AnsibleHostsPluginConfig{
				Modules:           []string{"a", "b"},
				FqdnOutputName:    "instance_fqdns",
				RecordsOutputName: "instance_records",
				Groups: map[string]*Group{
					"build": {"^build"},
					"ops":   {"^ops"},
					"zero":  {"0"},
					"one":   {"1"},
				},
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

			u := &AnsibleHostsPlugin{
				Config: v.config,
			}

			c := &plugin.StubContext{
				GlobalContext: &context.Context{
					Delegate: &exec.Cmd{
						Path: "./terraform",
					},
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

			expected, err := ioutil.ReadFile("./reference")
			check(err)
			got, err := ioutil.ReadFile("./hosts")

			if !reflect.DeepEqual(got, expected) {
				t.Fatalf("Files do not match, expected \n[%s], got \n[%s]", expected, got)
			}

		})
	}

}
func TestAnsibleHostsPlugin_Config(t *testing.T) {

	plugins := map[string]config.WithConfig{
		"ansible-hosts": &AnsibleHostsPlugin{},
	}

	cases := map[string]struct {
		input    string
		plugins  map[string]config.WithConfig
		expected *config.Config
	}{
		"empty": {
			`
executable test {
    plugin ansible-hosts {
	}
}
`,
			plugins,
			&config.Config{
				Executables: map[string]*config.Executable{
					"test": {
						Plugins: map[string]interface{}{
							"ansible-hosts": &AnsibleHostsPluginConfig{},
						},
					},
				},
			}},
		"simple": {
			`
executable test {
    plugin ansible-hosts {
		modules = [ "ops", "build" ]
		fqdn_output_name = "instance_fqdns"
		records_output_name = "instance_records"
		group "group1" {
			pattern = "pattern1"
		}
		group "group2" {
			pattern = "pattern2"
		}
		group "group3" {
		}
	}
}
`,
			plugins,
			&config.Config{
				Executables: map[string]*config.Executable{
					"test": {
						Plugins: map[string]interface{}{
							"ansible-hosts": &AnsibleHostsPluginConfig{
								Modules:           []string{"ops", "build"},
								FqdnOutputName:    "instance_fqdns",
								RecordsOutputName: "instance_records",
								Groups: map[string]*Group{
									"group1": &Group{
										Pattern: "pattern1",
									},
									"group2": &Group{
										Pattern: "pattern2",
									},
									"group3": &Group{},
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
