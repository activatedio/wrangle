package ansiblehosts

import (
	"strings"
	"testing"

	"github.com/activatedio/wrangle/config"
	"github.com/activatedio/wrangle/plugin"
)

var _ (plugin.Plugin) = (*AnsibleHostsPlugin)(nil)

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
