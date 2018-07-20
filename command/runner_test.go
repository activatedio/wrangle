package command

import "github.com/activatedio/wrangle/plugin"

var _ (plugin.Context) = (*internalContext)(nil)

var _ (plugin.Plugin) = (*runnerPlugin)(nil)
