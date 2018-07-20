package ansiblehosts

import (
	"sync"

	"github.com/activatedio/wrangle/plugin"
)

type AnsibleHostsPluginConfig struct {
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
