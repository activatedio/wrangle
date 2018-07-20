package plugin

type Registry interface {
	Get(name string) (Plugin, bool)
	Set(name string, plugin Plugin)
}

type DefaultRegistry map[string]Plugin

func (self DefaultRegistry) Get(name string) (Plugin, bool) {
	p, ok := self[name]
	return p, ok
}

func (self DefaultRegistry) Set(name string, plugin Plugin) {
	self[name] = plugin
}

