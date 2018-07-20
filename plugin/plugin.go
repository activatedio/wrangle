package plugin

import "github.com/activatedio/wrangle/context"

type Context interface {
	GetGlobalContext() *context.Context
	Next() error
}

type Plugin interface {
	GetConfig() interface{}
	Filter(c Context) error
}
