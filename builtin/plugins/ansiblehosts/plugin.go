package ansiblehosts

import (
	"sync"

	"strings"

	"errors"

	"os/exec"

	"encoding/json"

	"fmt"
	"regexp"
	"sort"

	"os"

	"github.com/activatedio/wrangle/plugin"
)

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

	if err := c.Next(); err != nil {
		return err
	}

	p := c.GetGlobalContext().Delegate.Path

	if !strings.HasSuffix(p, "terraform") {
		return errors.New("Executable must be terraform")
	}

	conf := self.Config

	var entries []*entry

	for _, m := range conf.Modules {
		cmd := exec.Command(p, "output", "-module="+m, "-json")
		out, err := cmd.Output()
		if err != nil {
			continue
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

	if len(entries) == 0 {
		return nil
	}

	f, err := os.Create("./hosts")
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString("[all]\n")

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

	return nil
}
