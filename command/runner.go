package command

import (
	"github.com/activatedio/wrangle/context"
	"github.com/activatedio/wrangle/plugin"
	"github.com/golang-collections/collections/stack"
)

type Runner interface {
	Run() error
}

type internalContext struct {
	globalContext *context.Context
	plugins       *stack.Stack
}

func (self *internalContext) GetGlobalContext() *context.Context {
	return self.globalContext
}

func (self *internalContext) Next() error {
	return self.plugins.Pop().(plugin.Plugin).Filter(self)
}

func newInternalContext(globalContext *context.Context, plugins []plugin.Plugin) *internalContext {

	c := &internalContext{
		globalContext: globalContext,
		plugins:       stack.New(),
	}

	for i := len(plugins); i > 0; i-- {
		c.plugins.Push(plugins[i-1])
	}

	return c
}

type runnerPlugin struct {
}

func (self *runnerPlugin) GetConfig() interface{} {
	return nil
}

func (self *runnerPlugin) Filter(c plugin.Context) error {

	return c.GetGlobalContext().Delegate.Run()
}

type DefaultRunner func(runner plugin.Plugin) plugin.Context

func NewDefaultRunner(globalContext *context.Context, plugins []plugin.Plugin) DefaultRunner {

	return func(runner plugin.Plugin) plugin.Context {
		return newInternalContext(globalContext, append(plugins, runner))
	}
}

func (self DefaultRunner) Run() error {

	c := self(&runnerPlugin{})
	return c.Next()
}
