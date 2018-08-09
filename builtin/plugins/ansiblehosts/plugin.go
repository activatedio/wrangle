package ansiblehosts

import (
	"sync"

	"github.com/activatedio/wrangle/plugin"
)

var config2 = `

---
instance_modules:
  - core_import: core_import_west
    instances: ops_instances_west
  - core_import: core_import_west
    instances: ns_instances_west
groups:
  all:
  ops:
    pattern: 'ops.*'
  microservice-controller:
    pattern: '^ops[a-z]+[0-2].*'
  microservice-controller-west:
    pattern: '^opswest[0-2].*'
  microservice-client:
    pattern: 'ops.*3.*'
  west:
    pattern: '.*west.*'
  vault:
    pattern: 'vault.*'
  bind:
    pattern: 'ns.*'
  bind-west:
    pattern: 'nswest.*'

`

type Group struct {
	Pattern string
}

type AnsibleHostsPluginConfig struct {
	Modules           []string          `hcl:"modules"`
	FqdnOutputName    string            `hcl:"fqdn_output_name"`
	RecordsOutputName string            `hcl:"records_output_name"`
	Groups            map[string]*Group `hcl:"group"`
}

type AnsibleHostsPlugin struct {
	ConfigLock sync.Mutex
	Config     *AnsibleHostsPluginConfig
}

func (self *AnsibleHostsPlugin) GetConfig() interface{} {

	self.ConfigLock.Lock()
	defer self.ConfigLock.Unlock()

	if self.Config == nil {
		self.Config = &AnsibleHostsPluginConfig{}
	}

	return self.Config
}

func (self *AnsibleHostsPlugin) Filter(c plugin.Context) error {

	return nil
}
