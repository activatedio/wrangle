package template

import (
	"sync"

	"os"

	"io/ioutil"

	"strings"

	"html/template"

	"bufio"

	"github.com/activatedio/wrangle/plugin"
	"gopkg.in/yaml.v2"
)

type TemplatePluginConfig struct {
	DataFile string
}

type TemplatePlugin struct {
	ConfigLock sync.Mutex
	Config     *TemplatePluginConfig
}

func (self *TemplatePlugin) GetConfig() interface{} {

	self.ConfigLock.Lock()
	defer self.ConfigLock.Unlock()

	if self.Config == nil {
		self.Config = &TemplatePluginConfig{}
	}

	return self.Config
}

func (self *TemplatePlugin) Filter(c plugin.Context) error {

	files, err := ioutil.ReadDir(".")

	if err != nil {
		return err
	}

	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".tft") {

			n := f.Name()

			dest, err := os.Create(n[:len(n)-1])
			defer dest.Close()

			if err != nil {
				return err
			}

			t, err := template.ParseFiles(f.Name())

			if err != nil {
				return err
			}

			w := bufio.NewWriter(dest)

			d, err := self.getData()

			if err != nil {
				return err
			}

			err = t.Execute(w, d)

			w.Flush()

			if err != nil {
				return err
			}
		}
	}

	return c.Next()

}

func (self *TemplatePlugin) getData() (interface{}, error) {

	p := self.Config.DataFile

	if p == "" {
		return new(interface{}), nil
	}

	dat, err := ioutil.ReadFile(p)

	if err != nil {
		return nil, err
	}

	v := make(map[string]interface{})

	yaml.Unmarshal(dat, &v)

	return v, nil
}