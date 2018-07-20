package command

import (
	"testing"

	"reflect"

	"github.com/activatedio/wrangle/context"
	"github.com/activatedio/wrangle/plugin"
)

var _ (plugin.Context) = (*internalContext)(nil)

var _ (plugin.Plugin) = (*runnerPlugin)(nil)

var _ (plugin.Plugin) = (*stubPlugin)(nil)

type stubPlugin struct {
	t               *testing.T
	id              string
	idListFunc      func(id string)
	expectedContext plugin.Context
}

func (self *stubPlugin) GetConfig() interface{} {
	return nil
}

func (self *stubPlugin) Filter(c plugin.Context) error {
	if c != self.expectedContext {
		self.t.Fatal("Unexpected context")
	}
	self.idListFunc(self.id)
	return nil
}

func TestInternalContext_Next(t *testing.T) {

	var idList []string

	idListFunc := func(name string) {
		idList = append(idList, name)
	}

	plugins := []plugin.Plugin{
		&stubPlugin{
			t:          t,
			id:         "a",
			idListFunc: idListFunc,
		},
		&stubPlugin{
			t:          t,
			id:         "b",
			idListFunc: idListFunc,
		},
	}

	u := newInternalContext(nil, plugins)

	for _, v := range plugins {
		v.(*stubPlugin).expectedContext = u
	}

	for i := 0; i < 2; i++ {
		if err := u.Next(); err != nil {
			t.Fatalf("Unexpected error %s", err)
		}
	}

	idListWant := []string{"a", "b"}

	if !reflect.DeepEqual(idList, idListWant) {
		t.Fatalf("Wanted id list %s, got %s", idListWant, idList)
	}

	if u.plugins.Len() != 0 {
		t.Fatal("Expected plugins to be empty")
	}
}

func TestNewDefaultRunner(t *testing.T) {

	gcWant := &context.Context{}

	plugins := []plugin.Plugin{
		&stubPlugin{
			id: "a",
		},
		&stubPlugin{
			id: "b",
		},
	}

	runner := &stubPlugin{
		id: "c",
	}

	u := NewDefaultRunner(gcWant, plugins)

	result := u(runner)

	gcGot := result.GetGlobalContext()

	if !reflect.DeepEqual(gcWant, gcGot) {
		t.Fatalf("Global context, want %+v, got %+v", gcWant, gcGot)
	}

	s := result.(*internalContext).plugins

	allPlugins := append(plugins, runner)

	for i := 0; i < len(allPlugins); i++ {
		want := allPlugins[i]
		got := s.Pop()
		if !reflect.DeepEqual(got, allPlugins[i]) {
			t.Fatalf("Iter %d, wanted %+v, got %+v", i, want, got)
		}
	}
}

var _ (Runner) = (*DefaultRunner)(nil)

func TestDefaultRunner_Run(t *testing.T) {

	var u DefaultRunner

	c := &plugin.StubContext{}

	u = func(runner plugin.Plugin) plugin.Context {
		return c
	}

	if u.Run() != nil {
		t.Fatal("Unexpected error")
	}

	if c.NextCallCount != 1 {
		t.Fatal("Expected next to be called one time")
	}

}
