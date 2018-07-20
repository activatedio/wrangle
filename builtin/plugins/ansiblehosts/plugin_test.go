package ansiblehosts

import (
	"testing"

	"github.com/activatedio/wrangle/plugin"
)

var _ (plugin.Plugin) = (*AnsibleHostsPlugin)(nil)

func TestAnsibleHostsPlugin(t *testing.T) {

	plugin.TestPlugin(t, &AnsibleHostsPlugin{}, &AnsibleHostsPluginConfig{},
		func(c interface{}) {

		})
}
