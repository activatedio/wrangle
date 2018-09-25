package template

import (
	"sync"
	"text/template"

	"os"

	"io/ioutil"

	"bufio"

	"strings"

	"github.com/activatedio/wrangle/context"
	"github.com/activatedio/wrangle/plugin"
	"gopkg.in/yaml.v2"
)

type SuffixConfig struct {
	ToSuffix string `hcl:"to-suffix"`
}

type TemplatePluginConfig struct {
	DataFile string                   `hcl:"data-file"`
	Suffixes map[string]*SuffixConfig `hcl:"suffix"`
}

type TemplatePlugin struct {
	ConfigLock sync.Mutex
	Config     *TemplatePluginConfig
}

func (self *TemplatePlugin) GetConfig() interface{} {

	self.ConfigLock.Lock()
	defer self.ConfigLock.Unlock()

	if self.Config == nil {
		self.Config = &TemplatePluginConfig{
			Suffixes: map[string]*SuffixConfig{},
		}
	}

	return self.Config
}

func (self *TemplatePlugin) Filter(c plugin.Context) error {

	files, err := ioutil.ReadDir(".")

	if err != nil {
		return err
	}

	for _, f := range files {
		if suffix, suffixConfig, ok := self.getSuffix(f.Name()); !f.IsDir() && ok {

			n := f.Name()

			dest, err := os.Create(n[:len(n)-len(suffix)] + suffixConfig.ToSuffix)
			defer dest.Close()

			if err != nil {
				return err
			}

			t, err := template.New("template").Funcs(map[string]interface{}{
				"join":    join,
				"project": project,
			}).ParseFiles(f.Name())

			if err != nil {
				return err
			}

			w := bufio.NewWriter(dest)

			d, err := self.getData(c.GetGlobalContext())

			if err != nil {
				return err
			}

			err = t.ExecuteTemplate(w, f.Name(), d)

			w.Flush()

			if err != nil {
				return err
			}
		}
	}

	return c.Next()

}

func (self *TemplatePlugin) getData(context *context.Context) (interface{}, error) {

	v := make(map[string]interface{})

	p := self.Config.DataFile

	if p != "" {

		dat, err := ioutil.ReadFile(p)

		if err != nil {
			return nil, err
		}

		yaml.Unmarshal(dat, &v)
	}

	for k, v2 := range context.Variables {
		v[k] = v2
	}

	return v, nil
}

var defaultSuffixes = map[string]*SuffixConfig{
	".tmpl": &SuffixConfig{
		ToSuffix: "",
	},
}

func (self *TemplatePlugin) getSuffix(name string) (string, *SuffixConfig, bool) {

	suffixes := self.Config.Suffixes

	if len(suffixes) == 0 {
		suffixes = defaultSuffixes
	}

	for k, v := range suffixes {
		if strings.HasSuffix(name, k) {
			return k, v, true
		}
	}

	return "", nil, false
}
