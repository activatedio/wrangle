package plugin

import (
	"reflect"
	"testing"

	"github.com/activatedio/wrangle/context"
)

type StubPlugin struct {
	Config      interface{}
	FilterCalls []struct {
		Context Context
		Error   error
	}
}

var _ Plugin = (*StubPlugin)(nil)

func (self *StubPlugin) GetConfig() interface{} {
	return self.Config
}

func (self *StubPlugin) Filter(c Context) error {
	return nil
}

func TestPlugin(t *testing.T, u Plugin, expectedConfig interface{},
	modifyConfig func(c interface{})) {

	t.Run("GetConfig", func(t *testing.T) {

		got := u.GetConfig()

		if !reflect.DeepEqual(expectedConfig, got) {
			t.Fatal("Expected empty config")
		}

		modifyConfig(got)

		if !reflect.DeepEqual(u.GetConfig(), got) {
			t.Fatal("Expected value to be the same the second time")
		}

	})
}

var _ (Context) = (*StubContext)(nil)

type StubContext struct {
	GlobalContext *context.Context
	NextCallCount int
}

func (self *StubContext) Next() error {
	self.NextCallCount++
	return nil
}

func (self *StubContext) GetGlobalContext() *context.Context {
	return self.GlobalContext
}
