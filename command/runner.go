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

type runnerPlugin struct {
}

func (self *runnerPlugin) GetConfig() interface{} {
	return nil
}

func (self *runnerPlugin) Filter(c plugin.Context) error {

	return c.GetGlobalContext().Delegate.Run()
}

type DefaultRunner []plugin.Plugin

func (self DefaultRunner) Run() error {

	c := &internalContext{
		plugins: stack.New(),
	}

	for i := len(self); i > 0; i-- {
		c.plugins.Push(self[i-1])
	}

	c.plugins.Push(&runnerPlugin{})

	return c.Next()
}
