package ansiblehosts

import (
	"sync"

	"strings"

	"errors"

	"os/exec"

	"encoding/json"

	"os"

	"fmt"
	"regexp"
	"sort"

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

type entry struct {
	name string
	ip   string
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

	f, err := os.Create("./hosts")
	if err != nil {
		return err
	}
	defer f.Close()

	p := c.GetGlobalContext().Delegate.Path

	if !strings.HasSuffix(p, "terraform") {
		return errors.New("Executable must be terraform")
	}

	conf := self.Config

	var entries []*entry

	for _, m := range conf.Modules {
		cmd := exec.Command(p, "-module="+m, "-json")
		out, err := cmd.Output()
		if err != nil {
			return err
		}
		j := make(map[string]interface{})
		err = json.Unmarshal(out, &j)
		if err != nil {
			return err
		}

		if _, ok := j[conf.FqdnOutputName]; !ok {
			break
		}
		if _, ok := j[conf.RecordsOutputName]; !ok {
			break
		}

		fqdns := j[conf.FqdnOutputName].(map[string]interface{})["value"].([]interface{})
		records := j[conf.RecordsOutputName].(map[string]interface{})["value"].([]interface{})

		for i, fqdn := range fqdns {

			entries = append(entries, &entry{
				name: fqdn.(string),
				ip:   records[i].([]interface{})[0].(string),
			})
		}
	}

	if len(entries) > 0 {
		f.WriteString("[all]\n")
	}

	for _, entry := range entries {
		f.WriteString(fmt.Sprintf("%s ansible_host=%s\n", entry.name, entry.ip))
	}

	if conf.Groups != nil {

		var keys []string
		for k := range conf.Groups {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		for _, k := range keys {

			v := conf.Groups[k]

			f.WriteString(fmt.Sprintf("\n[%s]\n", k))

			for _, entry := range entries {
				match, err := regexp.Match(v.Pattern, []byte(entry.name))
				if err != nil {
					return err
				}
				if match {
					f.WriteString(entry.name + "\n")
				}
			}

		}
	}

	c.Next()

	return nil
}
